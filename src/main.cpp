#include <iostream>

#include "cli.hpp"
#include "command_runner.hpp"
#include "http_client.hpp"

int main(int argc, char** argv) {
  CLI cli("0.0.1");
  cli.run(argc, argv);
  CommandRunner cmdRunner;
  std::cout << cmdRunner.run("ls .") << std::endl;

  HttpClient http("www.example.com", "80");
  std::string response = http.get("/");
  std::cout << response << "\n";

  return 0;
}
