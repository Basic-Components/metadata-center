<template>
  <div class="components">
    <el-row type="flex" justify="center">
      <h2>组件列表</h2>
    </el-row>
    <el-row type="flex" justify="center">
      <el-table :data="component_list" style="width: 100%">
        <el-table-column label="id" width="180" align="center">
          <template v-slot:default="scope">
            <span style="margin-left: 10px">{{ scope.row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="name" width="180" align="center">
          <template v-slot:default="scope">
            <div slot="reference" class="name-wrapper">
              <el-tag size="medium">{{ scope.row.name }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="version" width="180" align="center">
          <template v-slot:default="scope">
            <div slot="reference" class="name-wrapper">
              <el-tag size="medium">{{ scope.row.version }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center">
          <template v-slot:default="scope">
            <el-button size="mini" @click="handleEdit(scope.$index, scope.row)">详细</el-button>
            <!-- <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button> -->
          </template>
        </el-table-column>
        <el-table-column type="expand">
          <template v-slot:default="scope">
            <el-form label-position="left" inline class="demo-table-expand">
              <el-form-item label="简介">
                <span>{{ scope.row.desc || "没有简介" }}</span>
              </el-form-item>
            </el-form>
          </template>
        </el-table-column>
        
      </el-table>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "components",
  data() {
    return {
      component_list: []
    };
  },
  methods: {
    handleEdit(index, row) {
      let id = row.id.toString();
      this.$router.push({
        name: "componentinfo",
        params: { id }
      });
    },
    // handleDelete(index, row) {
    //   this.$axios.get(`/component/`)({ heroId: row.id })
    // },
  },
  created: function() {
    this.createdPromise = this.$axios.get(`/component/`);
  },
  mounted: async function() {
    try {
      let response = await this.createdPromise;
      this.component_list = response.data;
    } catch (error) {
      alert("数据未找到");
    }
  }
};
</script>