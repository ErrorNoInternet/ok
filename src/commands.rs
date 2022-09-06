use crate::database::Database;
use chrono::{Datelike, TimeZone};
use colored::Colorize;
use std::ops::Index;

pub fn statistics_command(db: &Database) {
    let mut records = Vec::new();
    let mut keys = match db.keys() {
        Ok(keys) => keys,
        Err(error) => {
            println!("Uh oh! There was an error: {}", error);
            return;
        }
    };
    for _ in 0..3 {
        let mut highest: (String, u128) = (String::new(), 0);
        for key in &keys {
            if key.starts_with("day.") {
                match db.get(key.to_string()) {
                    Ok(value) => {
                        let counter = match value.parse() {
                            Ok(counter) => counter,
                            Err(error) => {
                                println!("Uh oh! There was an error: {}", error);
                                return;
                            }
                        };
                        if counter >= highest.1 {
                            highest = (key.to_string(), counter);
                        }
                    }
                    Err(error) => {
                        println!("Uh oh! There was an error: {}", error);
                        return;
                    }
                }
            }
        }
        if !highest.0.is_empty() {
            let month = match highest
                .0
                .split(".")
                .collect::<Vec<&str>>()
                .index(1)
                .clone()
                .parse()
            {
                Ok(month) => month,
                Err(error) => {
                    println!("Uh oh! There was an error: {}", error);
                    return;
                }
            };
            let day = match highest
                .0
                .split(".")
                .collect::<Vec<&str>>()
                .index(2)
                .clone()
                .parse()
            {
                Ok(day) => day,
                Err(error) => {
                    println!("Uh oh! There was an error: {}", error);
                    return;
                }
            };
            let record_time = chrono::Local.ymd(chrono::Local::now().year(), month, day);
            let mut label = "OKs";
            if highest.1 == 1 {
                label = "OK"
            }
            records.push(format!(
                "{} {} {}",
                record_time.format("%B %d:").to_string().bold(),
                highest.1,
                label
            ));
            keys.remove(keys.iter().position(|x| *x == highest.0).unwrap());
        }
    }
}

pub fn leaderboard_join_command(_db: &Database) {
    println!("not implemented");
}

pub fn leaderboard_leave_command(_db: &Database) {
    println!("not implemented");
}
