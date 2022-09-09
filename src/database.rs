use std::{fs, path::Path};

pub struct Database {
    pub path: String,
}

impl Database {
    pub fn open(path: String) -> Result<Self, String> {
        match fs::create_dir_all(&path) {
            Ok(_) => Ok(Database { path }),
            Err(error) => Err(error.to_string()),
        }
    }

    pub fn set(&self, key: String, value: String) -> Result<(), String> {
        if key.starts_with(".") {
            return Err(String::from("key cannot begin with a dot!"));
        }
        match fs::write(Path::new(self.path.as_str()).join(key), value) {
            Ok(_) => Ok(()),
            Err(error) => Err(error.to_string()),
        }
    }

    pub fn get(&self, key: String) -> Result<String, String> {
        if key.starts_with(".") {
            return Err(String::from("key cannot begin with a dot!"));
        }
        match fs::read(Path::new(self.path.as_str()).join(key)) {
            Ok(value) => Ok(std::str::from_utf8(&value).unwrap().to_string()),
            Err(error) => Err(error.to_string()),
        }
    }

    pub fn delete(&self, key: String) -> Result<(), String> {
        if key.starts_with(".") {
            return Err(String::from("key cannot begin with a dot!"));
        }
        match fs::remove_file(Path::new(self.path.as_str()).join(key)) {
            Ok(_) => Ok(()),
            Err(error) => Err(error.to_string()),
        }
    }

    pub fn keys(&self) -> Result<Vec<String>, String> {
        let iterator = match fs::read_dir(&self.path) {
            Ok(iterator) => iterator,
            Err(error) => return Err(error.to_string()),
        };
        let mut keys = Vec::new();
        for key in iterator {
            keys.push(key.unwrap().file_name().into_string().unwrap());
        }
        Ok(keys)
    }
}
