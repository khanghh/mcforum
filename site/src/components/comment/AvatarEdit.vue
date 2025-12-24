<template>
  <div class="avatar-edit">
    <div class="avatar-view" :style="{ backgroundImage: 'url(' + value + ')' }">
      <div class="upload-view" @click="pickImage">
        <Icon name="TablerCloudUpload" class="w-6 h-6 mb-1" />
        <span>Click to modify</span>
      </div>
    </div>

    <input ref="uploadImage" accept="image/*" type="file" @input="uploadAvatar">
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  methods: {
    pickImage() {
      const currentObj = this.$refs.uploadImage
      currentObj.dispatchEvent(new MouseEvent('click'))
    },
    async uploadAvatar(e) {
      const files = e.target.files
      if (files.length <= 0) {
        return
      }
      try {
        // 上传头像
        const file = files[0]
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await useHttpPostMultipart('/api/upload', formData)

        // 设置头像
        await useHttpPostForm('/api/users/update/avatar', {
          body: {
            avatar: ret.url,
          },
        })

        // 重新加载数据
        const userStore = useUserStore()
        userStore.getCurrent()
        useMsgSuccess('头像更新成功')
      }
      catch (e) {
        console.error(e)
        useMsgError('头像更新失败')
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.avatar-edit {
  .avatar-view {
    width: 120px;
    height: 120px;
    background-size: cover;
    background-color: var(--bg-color2);
    border-radius: 50%;
    position: relative;

    &:hover {
      .upload-view {
        visibility: visible;
      }
    }

    .upload-view {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      color: var(--text-color);
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      border-radius: 50%;
      background-color: var(--bg-color-alpha);
      visibility: hidden;
      cursor: pointer;

      span {
        font-size: 13px;
        font-weight: 500;
      }
    }
  }

  input[type="file"] {
    display: none;
  }
}
</style>
