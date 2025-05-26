use std::io::{self, Write};

fn main() {
    loop {
        print!("> ");
        io::stdout().flush().unwrap();

        let mut input = String::new();
        io::stdin().read_line(&mut input).unwrap();
        let input = input.trim();

        let parts: Vec<&str> = input.split_whitespace().collect();

        match parts.as_slice() {
            ["quit"] => {
                println!("Goodbye!");
                break;
            }
            ["echo", rest @ ..] => {
                println!("{}", rest.join(" "));
            }
            ["add", nums @ ..] => {
                let sum: i32 = nums.iter().filter_map(|n| n.parse::<i32>().ok()).sum();
                println!("Sum: {}", sum);
            }
            _ => println!("Unknown command"),
        }
    }
}
