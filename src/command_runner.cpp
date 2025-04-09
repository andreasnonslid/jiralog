#include <boost/process.hpp>
#include <sstream>

#include "command_runner.h"

CommandRunner::CommandRunner(std::string shell,
                             std::vector<std::string> shell_args)
    : shell_(std::move(shell)), shell_args_(std::move(shell_args)) {}

void CommandRunner::set_shell(const std::string& shell,
                              const std::vector<std::string>& args) {
  shell_ = shell;
  shell_args_ = args;
}

std::string CommandRunner::run(const std::string& command) const {
  boost::process::ipstream pipe;
  std::ostringstream output;
  std::vector<std::string> args = shell_args_;
  args.push_back(command);

  boost::process::child process(boost::process::search_path(shell_),
                                boost::process::args(args),
                                boost::process::std_out > pipe);

  std::string line;
  while (pipe && std::getline(pipe, line)) {
    output << line << '\n';
  }

  process.wait();
  return output.str();
}
