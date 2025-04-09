#pragma once

#include <string>
#include <vector>

class CommandRunner {
public:
  CommandRunner(std::string shell = "bash", std::vector<std::string> shell_args = {"-c"});

  void set_shell(const std::string& shell, const std::vector<std::string>& args = {"-c"});
  std::string run(const std::string& command) const;

private:
  std::string shell_;
  std::vector<std::string> shell_args_;
};
