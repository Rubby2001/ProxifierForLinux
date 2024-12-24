<template>
  <div class="rules-app">
    <el-card class="card">
      <div class="header">
        <el-button type="primary" size="default" class="centered-button" @click="showRules"><el-icon><Setting /></el-icon> Rules</el-button>
        <el-select v-model="selectedRule" placeholder="Select Rule" size="default" class="rule-select" @visible-change="handleVisibleChange">
          <el-option v-for="rule in rules" :key="rule" :label="rule" :value="rule" />
        </el-select>
      </div>

      <el-divider />

      <el-form label-position="top" class="ip-input-form">
        <el-form-item label="Proxy IPs (one per line)">
          <el-input
              v-model="proxyIPs"
              type="textarea"
              placeholder="Enter proxy IPs here..."
              rows="4"
          />
        </el-form-item>
      </el-form>

      <div class="actions">
        <el-button type="success" @click="startProxy" size="default"><el-icon><Link /></el-icon>Start Proxy</el-button>
        <el-button type="danger" @click="cleanProxy" size="default"><el-icon><Link /></el-icon>Clean Proxies</el-button>
      </div>
    </el-card>
    <el-dialog v-model="showRuleDialog" title="Manage Rules" width="80%">
      <el-table :data="rulesList" style="width: 100%">
        <el-table-column prop="type" label="Type" width="80"/>
        <el-table-column prop="ip" label="IP" width="140"/>
        <el-table-column prop="port" label="Port" width="80"/>
        <el-table-column fixed="right" labl="Actions" width="140">
          <template #default="scope">
<!--            <el-button @click="editRule(scope.row)" type="success" size="small">Edit</el-button>-->
            <el-button @click="deleteRule(scope.row)" type="danger" size="small">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addRule">Add</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="showAddRuleDialog" :title="'Add Rule'" width="80%">
      <el-form :model="currentRule" label-width="100px">
        <el-form-item label="Type">
          <el-select v-model="currentRule.type" placeholder="Select Type">
            <el-option label="SOCKS5" value="socks5" />
            <el-option label="SOCKS4" value="socks4" />
            <el-option label="HTTP" value="http" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model="currentRule.ip" />
        </el-form-item>
        <el-form-item label="Port">
          <el-input v-model="currentRule.port" />
        </el-form-item>
        <el-form-item label="Need Auth?">
          <el-checkbox v-model="currentRule.needAuth">Need Auth</el-checkbox>
        </el-form-item>
        <el-form-item v-if="currentRule.needAuth" label="Username">
          <el-input v-model="currentRule.username" />
        </el-form-item>
        <el-form-item v-if="currentRule.needAuth" label="Password">
          <el-input v-model="currentRule.password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddRuleDialog = false">Cancel</el-button>
          <el-button type="primary" @click="saveRule">Save</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { Setting,Link } from '@element-plus/icons-vue';
import {GetRules, GetRuleList, SaveRules, DeleteRule, StartProxy, CleanProxy} from '../wailsjs/go/main/App'
import { ref } from 'vue';
import {ElMessage, ElMessageBox} from "element-plus";
const showRuleDialog = ref(false);
const showAddRuleDialog = ref(false);
const editingRule = ref(false);
const rules = ref([]);
const selectedRule = ref("");
const proxyIPs = ref("10.0.0.0/8\n172.16.0.0/12\n192.168.0.0/16\n");
const currentRule = ref({
  type: '',
  ip: '',
  port: '',
  needAuth: false,
  username: '',
  password: ''
});
const rulesList = ref([
  // 初始化一些规则数据
]);

const startProxy = async() => {
  const result = await StartProxy(selectedRule.value,proxyIPs.value)
  if (result === 1){
    ElMessage.success("Start proxy successfully.");
  }else{
    ElMessage.error("Start proxy error.");
  }
};

const cleanProxy = async () => {
  const result = await CleanProxy();
  if (result === 1){
    ElMessage.success("Clean proxy successfully.");
  }else{
    ElMessage.error("Clean proxy error.");
  }
};

const addRule = () => {
  showAddRuleDialog.value = true;
};
// const editRule = (rule) => {
//   // 复制要编辑的规则到 currentRule
//   currentRule.type = rule.type;
//   currentRule.ip = rule.ip;
//   currentRule.port = rule.port;
//   if (rule.username){
//     currentRule.needAuth = true;
//     currentRule.username = rule.username;
//     currentRule.password = rule.password;
//   }
//   editingRule.value = rule;
//   showAddRuleDialog.value = true;
// };

const deleteRule = (rule) => {
  ElMessageBox.confirm(`Are you sure to delete?`, 'Delete Confirm', {
    confirmButtonText: 'Yes',
    cancelButtonText: 'No',
    type: 'warning',
  }).then(async () => {
      const result = await DeleteRule(rule);
      if (result === 1){
        ElMessage.success("Delete successfully.");
      }else{
        ElMessage.error("Delete error");
      }
    rulesList.value = await GetRuleList();
  })
};

const saveRule = async() => {
  // // 保存或更新规则逻辑
  // if (editingRule.value) {
  //   if (!isValidIP(currentRule.value.ip)){
  //     ElMessage.error("IP format is incorrect")
  //     return
  //   }
  //   if (!isValidPort(currentRule.value.port)){
  //     ElMessage.error("The port should be an integer between 1-65535")
  //     return
  //   }
  //   const result = await SaveEditedRules(currentRule.value);
  //   if (result === 1){
  //     ElMessage.success("Successfully saved!")
  //   }else if(result === 2){
  //     ElMessage.error("The rule with this ip and port already exists!")
  //     return
  //   }else if(result ===5){
  //     ElMessage.error("Unknown error!")
  //     return
  //   }
  // } else {
  //   if (!isValidIP(currentRule.value.ip)){
  //     ElMessage.error("IP format is incorrect")
  //     return
  //   }
  //   if (!isValidPort(currentRule.value.port)){
  //     ElMessage.error("The port should be an integer between 1-65535")
  //     return
  //   }
  //   const result = await SaveRules(currentRule.value);
  //   if (result === 1){
  //     ElMessage.success("Successfully saved!")
  //   }else if(result === 2){
  //     ElMessage.error("The rule with this ip and port already exists!")
  //     return
  //   }else if(result ===5){
  //     ElMessage.error("Unknown error!")
  //     return
  //   }
  // }
  if (!isValidIP(currentRule.value.ip)){
    ElMessage.error("IP format is incorrect")
    return
  }
  if (!isValidPort(currentRule.value.port)){
    ElMessage.error("The port should be an integer between 1-65535")
    return
  }
  const result = await SaveRules(currentRule.value);
  if (result === 1){
    ElMessage.success("Successfully saved!")
  }else if(result === 2){
    ElMessage.error("The rule with this ip and port already exists!")
    return
  }else if(result ===5){
    ElMessage.error("Unknown error!")
    return
  }
  rulesList.value = await GetRuleList();
  showAddRuleDialog.value = false;
};
const handleVisibleChange = async() =>{
  rules.value = await GetRules();
}
const showRules = async()=>{
  rulesList.value = await GetRuleList();
  showRuleDialog.value = true
}
function isValidIPv4(ip) {
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  return ipv4Regex.test(ip);
}
function isValidIPv6(ip) {
  const ipv6Regex = /^(([0-9a-fA-F]{1,4}:){7}([0-9a-fA-F]{1,4}|:)|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:))$/;
  return ipv6Regex.test(ip);
}
function isValidIP(ip) {
  return isValidIPv4(ip) || isValidIPv6(ip);
}
function isValidPort(port) {
  // 使用正则表达式确保输入只包含数字
  const portRegex = /^(?:[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/;

  return portRegex.test(port);
}

</script>

<style scoped>
.rules-app {
  display: flex ;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  background-color: #f5f5f5;
  padding: 20px;
}
.card {
  width: 500px;
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.rule-select {
  width: 300px;
}
.ip-input-form {
  margin-top: 20px;
}
.actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}
.centered-button {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
