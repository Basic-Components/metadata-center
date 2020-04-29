const state = {
    current_index: "/"
}


// getters
const getters = {

}

// actions 定义业务逻辑
const actions = {
    changeCurrrentIndex: function (context, payload) {
        sessionStorage.setItem("current_index", payload.current_index)
    },
    loadCurrrentIndex: function (context) {
        let current_index = sessionStorage.getItem("current_index")
        if (current_index) {
            let payload = { current_index }
            context.commit('changeCurrrentIndex', payload)
        } else {
        }
    }
}

// mutations 定义数据状态的操作
const mutations = {
    changeCurrrentIndex(state, payload) {
        let current_index = payload.current_index
        state.current_index = current_index
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
