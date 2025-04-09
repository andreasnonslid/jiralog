#pragma once

#include <iostream>
#include <string>
#include <vector>

class CLI {
public:
    explicit CLI(std::string version);

    int run(int argc, char** argv);
    int dispatch(const std::vector<std::string>& args);

protected:
    virtual void handle_command(const std::string& cmd, const std::vector<std::string>& args);
    virtual void show_help();
    virtual void show_version();

    std::string version_;

private:
    bool repl_mode_ = false;
};
