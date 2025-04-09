#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <sstream>

#include "http_client.hpp"

namespace beast = boost::beast;
namespace http = beast::http;
namespace net = boost::asio;
using tcp = net::ip::tcp;

HttpClient::HttpClient(std::string host, std::string port)
    : host_(std::move(host)), port_(std::move(port)) {}

std::string HttpClient::get(const std::string& path) {
  try {
    net::io_context ioc;
    tcp::resolver resolver(ioc);
    beast::tcp_stream stream(ioc);

    auto const results = resolver.resolve(host_, port_);
    stream.connect(results);

    http::request<http::string_body> req{http::verb::get, path, 11};
    req.set(http::field::host, host_);
    req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

    http::write(stream, req);

    beast::flat_buffer buffer;
    http::response<http::string_body> res;
    http::read(stream, buffer, res);

    stream.socket().shutdown(tcp::socket::shutdown_both);

    return res.body();
  } catch (const std::exception& e) {
    return std::string("HTTP GET error: ") + e.what();
  }
}

std::string HttpClient::post(const std::string& path, const std::string& body,
                             const std::string& content_type) {
  try {
    net::io_context ioc;
    tcp::resolver resolver(ioc);
    beast::tcp_stream stream(ioc);

    auto const results = resolver.resolve(host_, port_);
    stream.connect(results);

    http::request<http::string_body> req{http::verb::post, path, 11};
    req.set(http::field::host, host_);
    req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
    req.set(http::field::content_type, content_type);
    req.body() = body;
    req.prepare_payload();

    http::write(stream, req);

    beast::flat_buffer buffer;
    http::response<http::string_body> res;
    http::read(stream, buffer, res);

    stream.socket().shutdown(tcp::socket::shutdown_both);

    return res.body();
  } catch (const std::exception& e) {
    return std::string("HTTP POST error: ") + e.what();
  }
}
