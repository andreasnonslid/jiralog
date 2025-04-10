#pragma once

#include <string>
#include <map>
#include <boost/beast/http.hpp>
#include <boost/asio/ssl.hpp>
#include <boost/beast/core.hpp>

struct HttpResponse {
    int returnCode;
    std::string body;
};

class HttpClient {
 public:
  HttpClient(const std::string& host, const std::string& port = "80", bool use_ssl = false);

  HttpResponse get(const std::string& path);
  HttpResponse post(const std::string& path, const std::string& body);

  HttpResponse sendRequestWithHeaders(boost::beast::http::verb method,
                                      const std::string& path,
                                      const std::string& body,
                                      const std::map<std::string, std::string>& headers);

 private:
  const std::string host_;
  const std::string port_;
  const bool use_ssl_;

  using Request = boost::beast::http::request<boost::beast::http::string_body>;
  using Response = boost::beast::http::response<boost::beast::http::string_body>;
  using TcpStream = boost::beast::tcp_stream;
  using SslStream = boost::asio::ssl::stream<TcpStream>;

  Request createRequest(boost::beast::http::verb method, const std::string& path, const std::string& body = "");
  HttpResponse sendRequest(boost::beast::http::verb method, const std::string& path, const std::string& body = "");
  void sslHandshake(SslStream& ssl_stream);
};
