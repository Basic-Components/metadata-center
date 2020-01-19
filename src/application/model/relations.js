import connection from "./core"
export default function () {
    let ServiceModel = connection.get_table("Service")
    let SchemaModel = connection.get_table("Schema")
    ServiceModel.hasMany(SchemaModel,{as: 'Schemas'})
    SchemaModel.belongsTo(ServiceModel)
}
