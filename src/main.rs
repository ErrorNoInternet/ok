mod commands;
mod database;

use clap::Command;
use colored::Colorize;
use database::Database;
use rand::Rng;

fn main() {
    let mut command = Command::new("ok").author("ErrorNoInternet").subcommand(
        Command::new("statistics")
            .about("See your OK statistics")
            .alias("stats")
            .alias("status"),
    );
    if cfg!(feature = "online") {
        command = command
            .clone()
            .subcommand(
                Command::new("join")
                    .about("Join the OK leaderboard")
                    .alias("submit"),
            )
            .subcommand(Command::new("leave").about("Leave the OK leaderboard"));
    }

    match command.get_matches().subcommand() {
        Some(("statistics", _)) => {
            commands::statistics_command(&load_database());
        }
        Some(("join", _)) => {
            commands::leaderboard_join_command(&load_database());
        }
        Some(("leave", _)) => {
            commands::leaderboard_leave_command(&load_database());
        }
        _ => ok(&load_database()),
    }
}

fn ok(db: &Database) {
    let counter: u128 = match db.get(String::from("counter")) {
        Ok(counter) => match counter.parse() {
            Ok(counter) => counter,
            Err(_) => 0,
        },
        Err(_) => {
            println!("Welcome, my friend, the the land of OKs. Here's your first OK:");
            0
        }
    };
    match db.set(String::from("counter"), (counter + 1).to_string()) {
        Ok(_) => (),
        Err(error) => println!("Uh oh! There was an error: {}", error),
    };
    print_rainbow("ok");
}

fn print_rainbow(text: &str) {
    let mut generator = rand::thread_rng();
    for letter in text.chars() {
        print!(
            "{}",
            letter.to_string().truecolor(
                generator.gen_range(100..=255),
                generator.gen_range(100..=255),
                generator.gen_range(100..=255)
            )
        )
    }
    println!();
}

fn load_database() -> Database {
    let mut database_path = String::from(".ok");
    if cfg!(windows) {
        database_path = format!("C:\\Users\\{}\\AppData\\Roaming\\ok", whoami::username())
    } else if cfg!(unix) {
        database_path = format!("/home/{}/.config/ok", whoami::username())
    }
    Database::open(database_path).unwrap()
}
