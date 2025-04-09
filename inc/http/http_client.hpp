class HttpClient {
 public:
  HttpClient(std::string host, std::string port = "80");

  std::string get(const std::string& path);
  std::string post(const std::string& path, const std::string& body,
                   const std::string& content_type = "application/json");

 private:
  std::string host_;
  std::string port_;
};
