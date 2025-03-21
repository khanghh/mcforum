<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item field="title" label="Title">
        <a-input v-model="form.title" />
      </a-form-item>
      <a-form-item field="url" label="Url">
        <a-input v-model="form.url" />
      </a-form-item>
      <a-form-item field="logo" label="Logo">
        <image-upload v-model="form.logo" />
      </a-form-item>
      <a-form-item field="summary" label="Summary">
        <a-textarea v-model="form.summary" allow-clear />
      </a-form-item>
      <a-form-item field="status" label="Status">
        <a-select v-model="form.status">
          <a-option :value="0">Enabled</a-option>
          <a-option :value="1">Disabled</a-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import ImageUpload from '@/components/ImageUpload.vue';

  const emit = defineEmits(['ok']);

  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const form = ref({
    id: '',
    title: '',
    url: '',
    logo: '',
    summary: '',
    status: 0,
  });
  const rules = {
    title: [{ required: true, message: 'Please fill in the title' }],
    url: [{ required: true, message: 'Please fill in the link' }],
  };

  const show = () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = 'Added';
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = 'Edit';

    try {
      form.value = await axios.get(`/api/admin/link/${id}`);
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
        ? '/api/admin/link/create'
        : '/api/admin/link/update';
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
