mod database;
use database::Database;

fn main() {}

fn load_database() -> Database {
    let mut database_path = String::from(".ok");
    if cfg!(windows) {
        database_path = format!("C:\\Users\\{}\\AppData\\Roaming\\ok", whoami::username())
    } else if cfg!(unix) {
        database_path = format!("/home/{}/.config/ok", whoami::username())
    }
    Database::open(database_path).unwrap()
}
