
<template>
  <div class="component-new">
    <el-row type="flex" justify="center">
      <h2>创建组件</h2>
    </el-row>
    <el-row type="flex" justify="center" v-if="alertVisible">
      <el-col :span="20">
        <el-alert
          :title="alert_title"
          :type="alert_type"
          :description="alert_desc"
          @close="closeAlert"
          center
        ></el-alert>
      </el-col>
    </el-row>
    <el-row type="flex" justify="center">
      <el-col :span="24">
        <el-form
          label-position="top"
          ref="component_info"
          :model="component_info"
          class="demo-table-expand"
        >
          <el-form-item
            label="组件名"
            prop="name"
            :rule="[
            { required: true, message: '组件名不能为空'},
            ]"
          >
            <el-input type="name" v-model="component_info.name"></el-input>
          </el-form-item>
          <el-form-item
            label="组件版本"
            prop="version"
            :rule="[
            { required: true, message: '组件版本不能为空'},
            ]"
          >
            <el-input type="version" v-model="component_info.version"></el-input>
          </el-form-item>
          <el-form-item label="简介">
            <el-input v-model="component_info.desc"></el-input>
          </el-form-item>
          <el-form-item label="组件镜像">
            <el-input v-model="component_info.image"></el-input>
          </el-form-item>
          <el-form-item
            label="组件模式"
            prop="schema"
            :rule="[
            { required: true, message: '组件模式不能为空'},
            ]"
          >
            <div>
              <p>选择模式文件上传</p>
              <input
                type="file"
                class="upload"
                @change="loadFile"
                ref="inputer"
                accept="application/json"
              />
              <el-input type="textarea" :autosize="true" v-model="component_info.schema"></el-input>
              <!-- <highlight-code lang="json" :code="schemaString"></highlight-code> -->
              <!-- <input
              type="file"
              class="upload"
              @change="loadFile"
              ref="inputer"
              accept="application/json"
            />
              <el-button v-if="checked" type="primary" icon="el-icon-folder-add" @click="cacheConfig"></el-button>-->
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">立即创建</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Ajv from "ajv";
const ajv = new Ajv();
export default {
  name: "componentnew",
  data() {
    return {
      alertVisible: false,
      alert_title: "",
      alert_type: "success",
      alert_desc: "",
      component_info: {
        name: "",
        version: "",
        schema: "",
        desc: "",
        image: ""
      }
    };
  },
  methods: {
    schemaObjcet: function() {
      return this.component_info.schema
        ? JSON.parse(this.component_info.schema)
        : {};
    },
    loadFile() {
      let inputDOM = this.$refs.inputer;
      if (inputDOM.files && inputDOM.files.length !== 0) {
        let file = inputDOM.files[0];
        let reader = new FileReader();
        reader.onload = e => {
          this.component_info.schema = reader.result;
          this.alertVisible = true;
          this.alert_title = "文件读取成功";
          this.alert_type = "success";
          this.alert_desc = "";
        };
        reader.readAsText(file);
      } else {
        this.alertVisible = true;
        this.alert_title = "文件读取失败";
        this.alert_type = "error";
        this.alert_desc = "";
      }
    },
    closeAlert() {
      this.alertVisible = false;
      this.alert_title = "";
      this.alert_type = "success";
      this.alert_desc = "";
    },
    async onSubmit() {
      if (
        !(
          this.component_info.name &&
          this.component_info.version &&
          this.component_info.schema
        )
      ) {
        this.alertVisible = true;
        this.alert_title = "上传组件信息失败";
        this.alert_type = "error";
        this.alert_desc = "组件信息不满足要求";
      } else {
        let query = {
          name: this.component_info.name,
          version: this.component_info.version,
          schema: this.schemaObjcet()
        };
        if (this.component_info.desc) {
          Object.assign(query, { desc: this.component_info.desc });
        }
        if (this.component_info.image) {
          Object.assign(query, { image: this.component_info.image });
        }
        let res = await this.$axios.post(
          `/component/`,
          JSON.stringify(query),
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        );
        if (res.status !== 200) {
          this.alertVisible = true;
          this.alert_title = "上传组件信息失败";
          this.alert_type = "error";
          this.alert_desc = res.data.error;
        } else {
          this.alertVisible = true;
          this.alert_title = "上传组件信息成功";
          this.alert_type = "success";
          this.alert_desc = "";
        }
      }
    }
  }
};
</script>