import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const defaultLoginResult = {
  token: null,
  user_id: null,
  user_name: null,
};

export default new Vuex.Store({
  state: {
    isLogin: false,
    loginResult: defaultLoginResult,
  },
  mutations: {
    init(state) {
      let loginResult = JSON.parse(localStorage.getItem("loginResult"));
      console.log(localStorage.getItem("loginResult"));
      if (loginResult != null) {
        state.loginResult = loginResult;
      }
    },
    login(state, loginResult) {
      state.loginResult = loginResult;
    },
    logout(state) {
      localStorage.removeItem("loginResult");
      state.loginResult = defaultLoginResult;
    },
  },
  actions: {},
  getters: {
    isLogin: (state) => state.loginResult.user_id !== null,
    userID: (state) => state.loginResult.user_id,
    username: (state) => state.loginResult.user_name,
    accessToken: (state) => state.loginResult.token,
  },
});
