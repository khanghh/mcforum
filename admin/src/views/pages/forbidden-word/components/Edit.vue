<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk">
    <a-form ref="formRef" :model="form" :rules="rules" auto-label-width>
      <a-form-item label="Type" field="type">
        <a-select v-model="form.type" placeholder="Type">
          <a-option label="Word" value="word" />
          <a-option label="Regex" value="regex" />
        </a-select>
      </a-form-item>

      <a-form-item label="Match" field="word">
        <a-input v-model="form.word" />
      </a-form-item>

      <a-form-item label="Remark" field="remark">
        <a-input v-model="form.remark" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
const emit = defineEmits(['ok']);

const appStore = useAppStore();
const formRef = ref();
const config = reactive({
  visible: false,
  isCreate: false,
  title: '',
});

const form = ref({
  type: undefined,
  word: undefined,
  remark: undefined,
});
const rules = {
  word: [{ required: true, message: 'Please enter a banned word' }],
};

const show = () => {
  formRef.value.resetFields();

  config.isCreate = true;
  config.title = 'Add prohibited words';
  config.visible = true;
};

const showEdit = async (id: any) => {
  formRef.value.resetFields();

  config.isCreate = false;
  config.title = 'Edit';

  try {
    form.value = await axios.get(`/api/admin/forbidden-word/${id}`);
  } catch (e: any) {
    useHandleError(e);
  }

  config.visible = true;
};

const handleCancel = () => {
  formRef.value.resetFields();
};
const handleBeforeOk = async (done: (closed: boolean) => void) => {
  const validateErr = await formRef.value.validate();
  if (validateErr) {
    done(false);
    return;
  }
  try {
    const url = config.isCreate
      ? '/api/admin/forbidden-word/create'
      : '/api/admin/forbidden-word/update';
    await axios.postForm<any>(url, jsonToFormData(form.value));
    useNotificationSuccess('Submit successfully');
    emit('ok');
    done(true);
  } catch (e: any) {
    useHandleError(e);
    done(false);
  }
};

defineExpose({
  show,
  showEdit,
});
</script>

<style lang="less" scoped></style>
