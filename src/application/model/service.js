import Sequelize from "sequelize"

const ServiceModel = {
    name: "Service",
    schema: {
        name: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '服务名'
        },
        description: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '服务的说明文字'
        },
        createdAt: Sequelize.DATE,
        updatedAt: Sequelize.DATE
    },
    meta: {
        tableName: 'service',
        comment: "服务数据模式管理",
        underscored: true
    }
}

export default ServiceModel