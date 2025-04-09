#include <sstream>

#include "cli.hpp"

CLI::CLI(std::string version) : version_(std::move(version)) {}

int CLI::run(int argc, char** argv) {
  repl_mode_ = (argc < 2);
  std::vector<std::string> args(argv + 1, argv + argc);
  return dispatch(args);
}

int CLI::dispatch(const std::vector<std::string>& args) {
  if (!args.empty()) {
    const std::string& command = args[0];
    std::vector<std::string> sub_args(args.begin() + 1, args.end());

    if (command == "exit" || command == "quit") {
      return 0;
    }

    if (command == "--help" || command == "help") {
      show_help();
    } else if (command == "--version" || command == "version") {
      show_version();
    } else {
      handle_command(command, sub_args);
    }
  }

  if (!repl_mode_) return 0;

  std::cout << "> ";
  std::string line;
  if (!std::getline(std::cin, line)) return 0;

  std::istringstream iss(line);
  std::vector<std::string> tokens;
  for (std::string token; iss >> token;) {
    tokens.push_back(token);
  }

  return dispatch(tokens);
}

void CLI::handle_command(const std::string& cmd,
                         const std::vector<std::string>& args) {
  std::cout << "Unknown command: " << cmd << "\n";
  show_help();
}

void CLI::show_help() {
  std::cout << "Available commands:\n"
            << "  help        Show this help message\n"
            << "  version     Show version info\n"
            << "\nUsage:\n"
            << "  <command> [args...]\n";
}

void CLI::show_version() { std::cout << "Version: " << version_ << "\n"; }
