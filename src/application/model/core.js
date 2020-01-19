import Sequelize from 'sequelize'

const defaultService = {
    name: "TEST",
    description: "一个测试"
}
const defaultSchemas = [
    {
        task: "card",
        version: "0.0.1",
        status: "dev",
        schema: {
            "description": "A representation of a person, company, organization, or place",
            "type": "object",
            "required": ["familyName", "givenName"],
            "properties": {
                "fn": {
                    "description": "Formatted Name",
                    "type": "string"
                },
                "familyName": {
                    "type": "string"
                },
                "givenName": {
                    "type": "string"
                },
                "additionalName": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "honorificPrefix": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "honorificSuffix": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "nickname": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "email": {
                    "type": "object",
                    "properties": {
                        "type": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                },
                "tel": {
                    "type": "object",
                    "properties": {
                        "type": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                },
                "adr": { "$ref": "http://example.com/address.schema.json" },
                "geo": { "$ref": "http://example.com/geographical-location.schema.json" },
                "tz": {
                    "type": "string"
                },
                "photo": {
                    "type": "string"
                },
                "logo": {
                    "type": "string"
                },
                "sound": {
                    "type": "string"
                },
                "bday": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "org": {
                    "type": "object",
                    "properties": {
                        "organizationName": {
                            "type": "string"
                        },
                        "organizationUnit": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    {
        task: "address",
        version: "0.0.1",
        status: "dev",
        schema: {
            "description": "An address similar to http://microformats.org/wiki/h-card",
            "type": "object",
            "properties": {
                "post-office-box": {
                    "type": "string"
                },
                "extended-address": {
                    "type": "string"
                },
                "street-address": {
                    "type": "string"
                },
                "locality": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "postal-code": {
                    "type": "string"
                },
                "country-name": {
                    "type": "string"
                }
            },
            "required": ["locality", "region", "country-name"],
            "dependencies": {
                "post-office-box": ["street-address"],
                "extended-address": ["street-address"]
            }
        }
    },
    {
        task: "geographical",
        version: "0.0.1",
        status: "dev",
        schema: {
            "title": "Longitude and Latitude Values",
            "description": "A geographical coordinate.",
            "required": ["latitude", "longitude"],
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number",
                    "minimum": -90,
                    "maximum": 90
                },
                "longitude": {
                    "type": "number",
                    "minimum": -180,
                    "maximum": 180
                }
            }
        }
    }
]



export class Connection {
    constructor() {
        this.TABLES = new Map()
        this.db = null
        this.callbacks = []
    }

    run_callback() {
        if (this.callbacks.length > 0) {
            for (let callback of this.callbacks) {
                callback(this.db)
            }
            this.callbacks = []
        }
    }

    init_url(url, options = {}) {
        this.db = new Sequelize(url, options)
        this.run_callback()
        return this.db
    }

    init_app(app, options = { logging: false }) {
        let dburl = app.config.get("DB_URL")
        if (dburl) {
            app.db = this.init_url(dburl, options)
            return app
        } else {
            throw "DB_URL not exist"
        }
    }

    add_callback(func) {
        this.callbacks.push(func)
    }

    register(Model) {
        const name = Model.name
        const schema = Model.schema
        const meta = Model.meta
        if (this.db) {
            this.TABLES.set(name, this.db.define(name, schema, meta))
        } else {
            let TABLES = this.TABLES
            this.add_callback(
                function (db) {
                    TABLES.set(name, db.define(name, schema, meta))
                }
            )
            if (this.db) {
                run_callback()
            }
        }
    }
    get_table(db_name) {
        return this.TABLES.get(db_name)
    }
    add_relation(relation_function) {
        if (this.db) {
            relation_function()
        } else {
            this.add_callback(function (db) {
                relation_function()
            })
            if (this.db) {
                run_callback()
            }
        }
    }
    async create_tables(table_name = null, safe = true) {
        if (safe) {
            if (table_name) {
                await this.TABLES.get(table_name).sync()
            } else {
                for (let [_, table] of this.TABLES.entries()) {
                    await table.sync()
                }
            }
        } else {
            if (table_name) {
                await this.TABLES.get(table_name).sync({
                    force: true
                })
            } else {
                for (let [_, table] of this.TABLES.entries()) {
                    await table.sync({
                        force: true
                    })
                }
            }
        }
    }
    async moke_data() {
        let Service = this.get_table("Service")
        let Schema = this.get_table("Schema")
        if (Service && Schema) {
            let rows = await Schema.findAll()
            if (rows.length === 0) {
                let service = await Service.create(defaultService)
                let schemas = []
                for (let defaultSchema of defaultSchemas){
                    let schema = await Schema.create(defaultSchema)
                    schemas.push(schema)
                }
                await service.setSchemas(schemas)
                console.log('{"msg":"table have moke data"}')
            } else {
                console.log('{"msg":"table already have data"}')
            }
        } else {
            throw "table user not registed"
        }
    }
    async init_connection(){
        await this.create_tables("Service")
        await this.create_tables("Schema")
        await this.moke_data()
    }
}
const connection = new Connection()
export default connection