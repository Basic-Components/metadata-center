import Sequelize from "sequelize"

const SchemaModel = {
    name: "Schema",
    schema: {
        task: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '数据模式针对的任务'
        },
        version: {
            type: Sequelize.STRING,
            allowNull: false,
            comment: '数据模式版本号'
        },
        env: {
            type: Sequelize.ARRAY(Sequelize.ENUM('dev', 'test', 'produce', 'release')),
            allowNull: false,
            comment: "数据模式服务的发布环境,分为dev,test,produce,release'"
        },
        direction:{
            type: Sequelize.ENUM('query', 'response', 'receive', 'send','ref'),
            comment: "数据模型的对应的方向,只有'query', 'response', 'receive', 'send','ref'"
        },
        labels:{
            type: Sequelize.ARRAY(Sequelize.STRING(100)),
            comment: "数据模型的标签,用于搜索"
        },
        description: {
            type: Sequelize.TEXT,
            allowNull: false,
            comment: '数据模式针对任务的描述'
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
        underscored: true,
        indexes: [
            {
                unique: true,
                fields: ['task','direction','version']
            }
        ]
    }
}

export default SchemaModel