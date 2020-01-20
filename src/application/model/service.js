import Sequelize from "sequelize"

const ServiceModel = {
    name: "Service",
    schema: {
        name: {
            type: Sequelize.STRING,
            primaryKey: true,
            comment: '服务名'
        },
        description: {
            type: Sequelize.TEXT,
            allowNull: false,
            comment: '服务的说明文字'
        },
        tags: {
            type: Sequelize.ARRAY(Sequelize.STRING(100)),
            comment: '服务的标签'
        },
        status: {
            type: Sequelize.ENUM('test','pre', 'online', 'deprecated'),
            allowNull: false,
            comment: "服务的状态,包括测试,预上线,上线和过期四种状态"
        },
        createdAt: Sequelize.DATE,
        updatedAt: Sequelize.DATE
    },
    meta: {
        tableName: 'service',
        comment: "服务数据,所谓服务指针对某一特定业务创建的一类程序.",
        underscored: true
    }
}

export default ServiceModel