use crate::database::Database;
use chrono::{Datelike, TimeZone};
use colored::Colorize;
use std::ops::Index;
use textplots::{Chart, Plot, Shape};

pub fn statistics_command(db: &Database) {
    let graph_history_days: u32 = 5;
    let graph_width = 100;
    let graph_height = 32;
    let graph_smoothness = 10;
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
    println!("{}", "OK Graph:".bold());
    let mut lines = chart.split("\n").collect::<Vec<&str>>();
    lines.remove(lines.len() - 1);
    for line in &lines {
        println!("  {}", line);
    }

    let mut graph_bottom_text = String::new();
    let mut current_time = chrono::Local::now() - chrono::Duration::days(graph_history_days as i64);
    for _ in 0..graph_history_days {
        current_time += chrono::Duration::days(1);
        let date = current_time.format(&date_format).to_string();
        graph_bottom_text.push_str(&date);
        graph_bottom_text.push_str(&" ".repeat(
            (((graph_width as f32 * 0.7 - (graph_width as f32 / 25.0)) as u32 / graph_history_days)
                - date.len() as u32) as usize
                - 1,
        ))
    }
    println!("  {}", graph_bottom_text.bold());
}

pub fn leaderboard_join_command(_db: &Database) {
    println!("not implemented");
}

pub fn leaderboard_leave_command(_db: &Database) {
    println!("not implemented");
}
