mod commands;
mod database;

use std::collections::HashMap;

use clap::Command;
use database::Database;

fn main() {
    let mut command = Command::new("ok")
        .author("ErrorNoInternet")
        .subcommand_required(true)
        .arg_required_else_help(true)
        .subcommand(
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
            commands::statistics_command();
        }
        Some(("join", _)) => {
            commands::leaderboard_join_command();
        }
        Some(("leave", _)) => {
            commands::leaderboard_leave_command();
        }
        _ => unreachable!(),
    }
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
