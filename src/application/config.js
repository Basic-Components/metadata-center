import fs from 'fs'
import path from 'path'
import {
    updateMap
} from "./utils"
export const DEFAULT_CONFIG = {
    "STATIC_PATH": "static",
    "PORT": 5000,
    "HOST": "0.0.0.0",
    "DB_URL": "postgresql://postgres:postgres@localhost:5432/test",
    "LOG_LEVEL": "warn",
    "ORIGIN": "http://localhost:8080"
}



export function init_config(options) {
    let config = new Map(Object.entries(DEFAULT_CONFIG))
    let defaultconfigPath = path.resolve("./", "config.json")
    if (fs.existsSync(defaultconfigPath)) {
        let temp_config = new Map(Object.entries(JSON.parse(fs.readFileSync(configPath))))
        updateMap(config, temp_config)
    } else {
        console.log("default config file not exist")
    }
    if (options.config) {
        if (fs.existsSync(options.config)) {
            let temp_config = new Map(Object.entries(JSON.parse(fs.readFileSync(options.config))))
            updateMap(config, temp_config)
        } else {
            throw "config file not exist"
        }
    }
    return config
}