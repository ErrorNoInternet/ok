use crate::database::Database;
use chrono::{Datelike, TimeZone};
use colored::Colorize;
use std::{io::Write, ops::Index};
use textplots::{Chart, Plot, Shape};

pub fn reset_command(db: &Database) {
    print!(
        "{} {} ",
        "Are you sure you want to reset your OK statistics?",
        "Y/N:".bold()
    );
    std::io::stdout().flush().unwrap();
    let mut input = String::new();
    std::io::stdin().read_line(&mut input).unwrap();
    if input.len() >= 1 {
        let letter = input.chars().nth(0).unwrap().to_lowercase().to_string();
        if letter == String::from("y") {
            print!(
                "{} {} ",
                "Are you very sure you want to reset your OK statistics?".red(),
                "Y/N:".bold().red()
            );
            std::io::stdout().flush().unwrap();
            let mut input = String::new();
            std::io::stdin().read_line(&mut input).unwrap();
            if input.len() >= 1 {
                let letter = input.chars().nth(0).unwrap().to_lowercase().to_string();
                if letter == String::from("y") {
                    let keys = match db.keys() {
                        Ok(keys) => keys,
                        Err(error) => {
                            println!("Uh oh! There was an error: {}", error);
                            return;
                        }
                    };
                    for key in keys {
                        match db.delete(key.clone()) {
                            Ok(_) => (),
                            Err(error) => println!("Unable to delete key ({}): {}", key, error),
                        }
                    }
                    println!("{}", "Your OK statistics have been reset!".bold());
                }
            }
        }
    }
}

pub fn statistics_command(db: &Database) {
    let graph_history_days: u32 = 5;
    let graph_width = 110;
    let graph_height = 32;
    let graph_smoothness = 20;
    let date_format = match std::env::var("OK_DATE") {
        Ok(date_format) => date_format,
        Err(_) => String::from("%b %d"),
    };

    let current_time = chrono::Local::now();
    let current_day_counter = match db.get(format!(
        "day.{}.{}",
        current_time.month(),
        current_time.day()
    )) {
        Ok(counter) => match counter.parse() {
            Ok(counter) => counter,
            Err(_) => 0,
        },
        Err(_) => 0,
    };
    println!(
        "{} {}",
        "OKs Today:".bold(),
        current_day_counter.to_string().blue()
    );
    let counter: i64 = match db.get(String::from("counter")) {
        Ok(counter) => match counter.parse() {
            Ok(counter) => counter,
            Err(_) => 0,
        },
        Err(_) => 0,
    };
    println!("{} {}", "OK Counter:".bold(), counter.to_string().blue());

    let mut keys = match db.keys() {
        Ok(keys) => keys,
        Err(error) => {
            println!("Uh oh! There was an error: {}", error);
            return;
        }
    };
    if keys.len() > 0 {
        println!("{}", "OK Records:".bold());
        for i in 0..3 {
            let mut highest: (String, i64) = (String::new(), 0);
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
                println!(
                    "  {}. {} {} {}",
                    (i + 1).to_string().bold().blue(),
                    (record_time.format(&date_format).to_string() + " -").bold(),
                    (highest.1).to_string().blue(),
                    label
                );
                keys.remove(keys.iter().position(|x| *x == highest.0).unwrap());
            }
        }
    }

    let mut current_time = chrono::Local::now() - chrono::Duration::days(graph_history_days as i64);
    let mut data = Vec::new();
    for _ in 0..graph_history_days - 1 {
        current_time += chrono::Duration::days(1);

        let current_day_counter: i64 = match db.get(format!(
            "day.{}.{}",
            current_time.month(),
            current_time.day()
        )) {
            Ok(current_day_counter) => match current_day_counter.parse() {
                Ok(current_day_counter) => current_day_counter,
                Err(_) => 0,
            },
            Err(_) => 0,
        };
        let next_day_key = format!("day.{}.{}", current_time.month(), current_time.day() + 1);
        let next_day_counter: i64 = match db.get(next_day_key.clone()) {
            Ok(next_day_counter) => match next_day_counter.parse() {
                Ok(next_day_counter) => next_day_counter,
                Err(_) => 0,
            },
            Err(_) => 0,
        };

        let difference = (next_day_counter - current_day_counter) as f32 / graph_smoothness as f32;
        for i in 0..graph_smoothness {
            data.push(current_day_counter as f32 + i as f32 * difference);
        }
    }
    let chart = Chart::new(graph_width, graph_height, 0.0, data.len() as f32)
        .lineplot(&Shape::Continuous(Box::new(|x| {
            *data.index(x as usize) as f32
        })))
        .to_string()
        .trim()
        .to_owned();
    let mut str_lines = chart.split("\n").collect::<Vec<&str>>();
    str_lines.remove(str_lines.len() - 1);
    let mut lines = Vec::new();
    for line in str_lines {
        lines.push(line.to_string());
    }

    println!("{}", "OK Graph:".bold());
    let mut index = 0;
    for line in lines.clone().iter_mut() {
        let first_character = match line.chars().collect::<Vec<char>>().iter().nth(0) {
            Some(character) => character.to_owned(),
            None => ' ',
        };
        match &line.strip_prefix(first_character) {
            Some(new_line) => *line = new_line.to_string(),
            None => (),
        };

        if index == 0 || index == lines.len() - 1 {
            let mut line_index = line.chars().count();
            let mut character = match line.chars().collect::<Vec<char>>().iter().nth(line_index) {
                Some(character) => character.to_owned(),
                None => ' ',
            };
            while character != '.' {
                line_index -= 1;
                character = match line.chars().collect::<Vec<char>>().iter().nth(line_index) {
                    Some(character) => character.to_owned(),
                    None => ' ',
                };
                match &line.strip_suffix(character) {
                    Some(new_line) => *line = new_line.to_string(),
                    None => (),
                }
            }
            while line.chars().count() < lines.index(1).chars().count() + 3 {
                line.insert(line.len() - 2, ' ')
            }
        } else {
            line.push_str("    ");
            line.replace_range(line.len() - 1..line.len(), "|");
        }

        println!("  |{}", line);
        index += 1;
    }

    let mut graph_bottom_text = String::new();
    let mut current_time = chrono::Local::now() - chrono::Duration::days(graph_history_days as i64);
    for _ in 0..graph_history_days {
        current_time += chrono::Duration::days(1);
        let date = current_time.format(&date_format).to_string();
        graph_bottom_text.push_str(&date);
        graph_bottom_text.push_str(&" ".repeat(
            (((graph_width as f32 * 0.7 - (graph_width as f32 / 20.0)) as u32 / graph_history_days)
                - date.len() as u32) as usize
                - 1,
        ))
    }
    println!("   {}", graph_bottom_text.bold());
}

pub fn leaderboard_join_command(_db: &Database) {
    println!("not implemented");
}

pub fn leaderboard_leave_command(_db: &Database) {
    println!("not implemented");
}
