import Sequelize from "sequelize"

const SchemaModel = {
    name: "Schema",
    schema: {
        version: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '数据模式版本号'
        },
        task: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '数据模式针对的任务'
        },
        tags:{
            type: Sequelize.ARRAY(Sequelize.STRING(100)),
            comment: '数据模型的标签'
        },
        description: {
            type: Sequelize.TEXT,
            allowNull: false,
            comment: '数据模式针对任务的描述'
        },
        status: {
            type: Sequelize.ENUM('dev', 'test', 'produce', 'release'),
            allowNull: false,
            comment: "数据模式的发布状态,分为dev,test,produce,release'"
        },
        schema: {
            type: Sequelize.JSONB,
            comment: '格式描述用的json schema数据模式文本'
        },
        createdAt: Sequelize.DATE,
        updatedAt: Sequelize.DATE
    },
    meta: {
        tableName: 'schema',
        comment: "服务数据模式管理",
        underscored: true
    }
}

export default SchemaModel