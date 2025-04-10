#include <iostream>
#include <map>

#include "cli.hpp"
#include "command_runner.hpp"
#include "http_client.hpp"

std::string base64Encode(const std::string &input);

int main(int argc, char **argv) {
  CLI cli("0.0.1");
  cli.run(argc, argv);
  CommandRunner cmdRunner;
  std::cout << cmdRunner.run("ls .") << std::endl;

  std::string jiraHost = "autostore.atlassian.net";
  std::string email = "andreas.havardsen@autostoresystem.com";
  std::string apiToken = "some token";
  std::string issueKey = "TIME-25";

  std::string credentials = email + ":" + apiToken;
  std::string authHeader = "Basic " + base64Encode(credentials);

  std::map<std::string, std::string> headers;
  headers["Authorization"] = authHeader;
  headers["Content-Type"] = "application/json";
  headers["Accept"] = "application/json";

  const std::string apiPath = "/rest/api/3/issue/" + issueKey;

  HttpClient jiraClient(jiraHost, "443", true);

  HttpResponse response = jiraClient.sendRequestWithHeaders(
      boost::beast::http::verb::get, apiPath, "", headers);

  if (response.returnCode == 0)
    std::cout << "Jira Issue Details:\n" << response.body << "\n";
  else
    std::cout << "Jira API Request Failed:\n" << response.body << "\n";

  return 0;
}

std::string base64Encode(const std::string &input) {
  static const std::string base64_chars =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
  std::string ret;
  int val = 0, valb = -6;
  for (unsigned char c : input) {
    val = (val << 8) + c;
    valb += 8;
    while (valb >= 0) {
      ret.push_back(base64_chars[(val >> valb) & 0x3F]);
      valb -= 6;
    }
  }
  if (valb > -6) ret.push_back(base64_chars[((val << 8) >> (valb + 8)) & 0x3F]);
  while (ret.size() % 4) ret.push_back('=');
  return ret;
}
