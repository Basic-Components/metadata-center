import Sequelize from 'sequelize'
import { selfService, selfSchemas } from './default_data'


const Op = Sequelize.Op;
const operatorsAliases = {
    $eq: Op.eq,
    $ne: Op.ne,
    $gte: Op.gte,
    $gt: Op.gt,
    $lte: Op.lte,
    $lt: Op.lt,
    $not: Op.not,
    $in: Op.in,
    $notIn: Op.notIn,
    $is: Op.is,
    $like: Op.like,
    $notLike: Op.notLike,
    $iLike: Op.iLike,
    $notILike: Op.notILike,
    $regexp: Op.regexp,
    $notRegexp: Op.notRegexp,
    $iRegexp: Op.iRegexp,
    $notIRegexp: Op.notIRegexp,
    $between: Op.between,
    $notBetween: Op.notBetween,
    $overlap: Op.overlap,
    $contains: Op.contains,
    $contained: Op.contained,
    $adjacent: Op.adjacent,
    $strictLeft: Op.strictLeft,
    $strictRight: Op.strictRight,
    $noExtendRight: Op.noExtendRight,
    $noExtendLeft: Op.noExtendLeft,
    $and: Op.and,
    $or: Op.or,
    $any: Op.any,
    $all: Op.all,
    $values: Op.values,
    $col: Op.col
};


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

    init_app(app, options = { logging: false, operatorsAliases }) {
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
                let service = await Service.create(selfService)
                let schemas = []
                for (let selfSchema of selfSchemas) {
                    let schema = await Schema.create(selfSchema)
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
    async init_connection() {
        await this.create_tables("Service")
        await this.create_tables("Schema")
        await this.moke_data()
    }
}
const connection = new Connection()
export default connection