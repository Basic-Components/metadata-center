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
            type: Sequelize.TEXT,
            allowNull: false,
            comment: '服务的说明文字'
        },
        tag: {
            type: Sequelize.STRING(100),
            comment: '服务的标记,用于区分服务的版本和状态,类比docker image的tag'
        },
        labels: {
            type: Sequelize.ARRAY(Sequelize.STRING(100)),
            comment: '服务的标签,用于搜索查询'
        },
        usage_state: {
            type: Sequelize.ENUM("developing", "using", "deprecated"),
            allowNull: false,
            comment: '服务的状态,用于搜索查询'
        },
        createdAt: Sequelize.DATE,
        updatedAt: Sequelize.DATE
    },
    meta: {
        tableName: 'service',
        comment: "服务数据,所谓服务指针对某一特定业务创建的一类程序.",
        underscored: true,
        indexes: [
            {
                unique: true,
                fields: ['name', 'tag']
            }
        ]
    }
}

export default ServiceModel