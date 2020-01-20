import connection from '../model'
import request from 'request'
import Ajv from 'ajv'
function loadSchema(uri) {
    return request.get(uri).then(function (res) {
        if (res.statusCode >= 400)
            throw new Error('Loading error: ' + res.statusCode);
        return res.body;
    });
}
const SchemaValidator = new Ajv({ loadSchema: loadSchema })

async function getFromTable(table, where, ...attributes) {
    console.log(where)
    let options
    if (attributes.length === 0) {
        if (Object.keys(where).length === 0) {
            options = {
                order: [['updatedAt', 'DESC']]
            }
        } else {
            options = {
                order: [['updatedAt', 'DESC']],
                where: where
            }
        }

    } else {
        if (Object.keys(where).length === 0) {
            options = {
                attributes,
                order: [['updatedAt', 'DESC']]
            }
        } else {
            options = {
                attributes,
                order: [['updatedAt', 'DESC']],
                where: where
            }
        }
    }
    let schemas = await connection.get_table(table).findAll(options)
    return schemas.map((item) => item.dataValues)
}
async function getFromTableByid(table, serviceId, ...attributes) {
    if (attributes.length === 0) {
        let service = await connection.get_table(table).findByPk(serviceId)
        return service.dataValues
    } else {
        let options = {
            attributes,
        }
        let service = await connection.get_table(table).findByPk(serviceId, options)
        if (service) {
            return service.dataValues
        } else {
            return null
        }

    }
}

async function getService(where, ...attributes) {
    return await getFromTable("Service", where, ...attributes)
}
async function getServiceByid(where, ...attributes) {
    return await getFromTableByid("Service", where, ...attributes)
}
async function getSchema(where, ...attributes) {
    return await getFromTable("Schema", where, ...attributes)
}



ajv.compileAsync(schema).then(function (validate) {
    var valid = validate(data);
    // ...
});

function loadSchema(uri) {
    return request.json(uri).then(function (res) {
        if (res.statusCode >= 400)
            throw new Error('Loading error: ' + res.statusCode);
        return res.body;
    });
}

export { getService, getServiceByid, getSchema, SchemaValidator }
