<template>
  <div class="component-info">
    <el-row type="flex" justify="center">
      <h2>组件详情</h2>
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
          <el-form-item label="组件名">
            <el-input v-model="component_info.name" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item label="组件版本">
            <el-input v-model="component_info.version" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item label="简介">
            <el-input v-model="component_info.desc"></el-input>
          </el-form-item>
          <el-form-item label="组件镜像">
            <el-input v-model="component_info.image"></el-input>
          </el-form-item>
          <el-form-item label="组件模式">
            <div>
              <p>选择配置文件进行校验</p>
              <input
                type="file"
                class="upload"
                @change="loadFile"
                ref="inputer"
                accept="application/json"
              />

              <highlight-code lang="json" :code="schemaString"></highlight-code>
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
            <el-button type="primary" @click="onSubmit">修改</el-button>
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
  name: "componentinfo",
  data() {
    return {
      component_info: {},
      alertVisible: false,
      alert_title: "",
      alert_type: "success",
      alert_desc: "",
      checked: false
    };
  },
  props: {
    id: {
      type: String,
      default: null
    }
  },
  computed: {
    // 计算属性的 getter
    schemaString: function() {
      return this.component_info.schema
        ? JSON.stringify(this.component_info.schema, null, 2)
        : "{}";
    }
  },
  methods: {
    loadFile() {
      let inputDOM = this.$refs.inputer;
      if (inputDOM.files && inputDOM.files.length !== 0) {
        let file = inputDOM.files[0];
        let reader = new FileReader();
        reader.onload = e => {
          let validater = ajv.compile(this.component_info.schema);
          let valid = validater(JSON.parse(reader.result));
          if (!valid) {
            this.alertVisible = true;
            this.alert_title = "验证失败";
            this.alert_type = "error";
            let error_infos = validater.errors.map(
              error_info => {
                let schemaPath = error_info.schemaPath;
                let error_msg = error_info.message;
                let temp = `Path: ${schemaPath}----Msg: ${error_msg}`;
                return temp;
              }
            );
            let alert_desc = error_infos.join("||");
            this.alert_desc = alert_desc;
            this.checked = false;
          } else {
            this.alertVisible = true;
            this.alert_title = "验证成功";
            this.alert_type = "success";
            this.alert_desc = "配置符合模式";
            this.checked = true;
          }
        };
        reader.readAsText(file);
      } else {
        this.alertVisible = true;
        this.alert_title = "文件读取失败";
        this.alert_type = "error";
        this.alert_desc = "";
        this.checked = false;
      }
    },
    closeAlert() {
      this.alertVisible = false;
      this.alert_title = "";
      this.alert_type = "success";
      this.alert_desc = "";
    },
    async onSubmit() {
      let Id = Number(this.id)
      if (
        (!this.component_info.desc) &&
        (!this.component_info.image)
      ) {
        this.alertVisible = true;
        this.alert_title = "更新组件信息失败";
        this.alert_type = "error";
        this.alert_desc = "组件信息不满足要求";
      } else {
        let query = {
        };
        if (this.component_info.desc) {
          Object.assign(query, { desc: this.component_info.desc });
        }
        if (this.component_info.image) {
          Object.assign(query, { image: this.component_info.image });
        }
        let res = await this.$axios.put(
          `/component/${Id}`,
          JSON.stringify(query),
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        );
        if (res.status !== 200) {
          this.alertVisible = true;
          this.alert_title = "更新组件信息失败";
          this.alert_type = "error";
          this.alert_desc = res.data.error;
        } else {
          this.alertVisible = true;
          this.alert_title = "更新组件信息成功";
          this.alert_type = "success";
          this.alert_desc = "";
        }
      }
    }
  },
  created: async function() {
    let Id = Number(this.id);
    let resp = await this.$axios.get(`/component/${Id}`);
    this.component_info = resp.data;
  }
};
</script>