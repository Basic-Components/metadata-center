<template>
  <div id="app">
    <el-container>
      <el-aside width="200px">
        <el-menu
          :default-active="activeIndex"
          class="el-menu-demo"
          mode="vertical"
          router
          @select="changeIndex"
        >
        <el-menu-item index="/">说明</el-menu-item>
        <el-submenu index="1">
          <template slot="title">
            <span>组件</span>
          </template>
          <el-menu-item index="/components">全部组件</el-menu-item>
          <el-menu-item index="/component-new">创建组件</el-menu-item>
        </el-submenu>
        <el-submenu index="2">
          <template slot="title">
            <span>服务</span>
          </template>
          <el-menu-item index="/services">全部服务</el-menu-item>
          <el-menu-item index="/service-new">创建服务</el-menu-item>
        </el-submenu>
        </el-menu>
      </el-aside>
      <el-main>
        <transition name="slide" mode="out-in" appear>
          <router-view></router-view>
        </transition>
      </el-main>
    </el-container>
  </div>
</template>

<script>
export default {
  name: "app",
  data() {
    return {
      activeIndex: "/"
    };
  },
  methods: {
    changeIndex: function(index, indexPath) {
      this.$store.dispatch("menu/changeCurrrentIndex", {
        current_index: index
      });
    }
  },
  created: function() {
    this.$store.dispatch("menu/loadCurrrentIndex");
    this.activeIndex = this.$store.state.menu.current_index;
  }
};
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
