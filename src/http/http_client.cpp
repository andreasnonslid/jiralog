#include <boost/asio/connect.hpp>
#include <boost/asio/deadline_timer.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/asio/ssl.hpp>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>

#include "http_client.hpp"

namespace beast = boost::beast;
namespace http = beast::http;
namespace net = boost::asio;

HttpClient::HttpClient(const std::string& host, const std::string& port,
                       bool use_ssl)
    : host_(host), port_(port), use_ssl_(use_ssl) {}

HttpClient::Request HttpClient::createRequest(boost::beast::http::verb method,
                                              const std::string& path,
                                              const std::string& body) {
  Request req{method, path, 11};
  req.set(http::field::host, host_);
  req.set(http::field::user_agent, "HttpClient/1.0");
  if (!body.empty()) {
    req.body() = body;
    req.prepare_payload();
  }
  return req;
}

void HttpClient::sslHandshake(SslStream& ssl_stream) {
  ssl_stream.handshake(boost::asio::ssl::stream_base::client);
}

HttpResponse HttpClient::sendRequest(boost::beast::http::verb method,
                                     const std::string& path,
                                     const std::string& body) {
  return sendRequestWithHeaders(method, path, body, {});
}

HttpResponse HttpClient::sendRequestWithHeaders(
    boost::beast::http::verb method, const std::string& path,
    const std::string& body,
    const std::map<std::string, std::string>& headers) {
  try {
    net::io_context ioc;
    net::ip::tcp::resolver resolver(ioc);
    TcpStream stream(ioc);
    std::shared_ptr<SslStream> ssl_stream;

    auto const results = resolver.resolve(host_, port_);
    stream.connect(results);

    Request req = createRequest(method, path, body);
    for (const auto& header : headers) req.set(header.first, header.second);

    if (use_ssl_) {
      boost::asio::ssl::context ctx(boost::asio::ssl::context::tlsv12_client);
      ctx.set_options(boost::asio::ssl::context::default_workarounds |
                      boost::asio::ssl::context::no_sslv2 |
                      boost::asio::ssl::context::no_sslv3);
      ssl_stream = std::make_shared<SslStream>(std::move(stream), ctx);

      if (!SSL_set_tlsext_host_name(ssl_stream->native_handle(),
                                    host_.c_str())) {
        boost::system::error_code ec{static_cast<int>(::ERR_get_error()),
                                     boost::asio::error::get_ssl_category()};
        throw boost::system::system_error{ec};
      }

      sslHandshake(*ssl_stream);
      http::write(*ssl_stream, req);
    } else {
      http::write(stream, req);
    }

    beast::flat_buffer buffer;
    Response res;
    if (use_ssl_)
      http::read(*ssl_stream, buffer, res);
    else
      http::read(stream, buffer, res);

    if (use_ssl_) {
      boost::system::error_code ec;
      ssl_stream->shutdown(ec);
    } else {
      stream.socket().shutdown(net::socket_base::shutdown_both);
    }
    return HttpResponse{0, res.body()};
  } catch (const std::exception& e) {
    return HttpResponse{1, std::string("Error: ") + e.what()};
  }
}

HttpResponse HttpClient::get(const std::string& path) {
  return sendRequest(http::verb::get, path);
}

HttpResponse HttpClient::post(const std::string& path,
                              const std::string& body) {
  return sendRequest(http::verb::post, path, body);
}
