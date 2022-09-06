use std::ops::Index;

use crate::database::Database;
use colored::Colorize;

pub fn statistics_command(db: &Database) {
    let mut records = Vec::new();
    let mut keys = match db.keys() {
        Ok(keys) => keys,
        Err(error) => {
            println!("Uh oh! There was an error: {}", error);
            return;
        }
    };
    for i in 0..3 {
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
            let month = highest.0.split(".").collect::<Vec<&str>>().index(1).clone();
            let day = highest.0.split(".").collect::<Vec<&str>>().index(2).clone();
            let mut label = "OKs";
            if highest.1 == 1 {
                label = "OK"
            }
            records.push(format!(
                "{} {} {}",
                format!("{}/{}:", month, day).bold(),
                highest.1,
                label
            ));
            keys.remove(keys.iter().position(|x| *x == highest.0).unwrap());
        }
    }
    for record in records {
        println!("{}", record);
    }
}

pub fn leaderboard_join_command(db: &Database) {
    println!("not implemented");
}

pub fn leaderboard_leave_command(db: &Database) {
    println!("not implemented");
}
