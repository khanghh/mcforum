<template>
  <ClientOnly>
    <el-dropdown v-if="hasPermission" trigger="click" @command="handleCommand">
      <span class="el-dropdown-link">{{ $t('publish.manage') }}</span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="edit">{{ $t('publish.action.edit') }}</el-dropdown-item>
          <el-dropdown-item command="delete">{{ $t('publish.action.delete') }}</el-dropdown-item>
          <el-dropdown-item command="pin">{{ $t('publish.action.pin') }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </ClientOnly>
</template>

<script setup>
const i18n = useI18n();
const props = defineProps({
  article: {
    type: Object,
    required: true,
  },
});

const userStore = useUserStore();
const isOwner = userIsOwner(userStore.user);
const isAdmin = userIsAdmin(userStore.user);
const isArticleOwner = computed(() => {
  return userStore.user && userStore.user.id === props.article.user.id;
});
const hasPermission = computed(() => {
  return isArticleOwner || isOwner || isAdmin;
});

async function handleCommand(command) {
  if (command === "edit") {
    editArticle();
  } else if (command === "delete") {
    deleteArticle();
  } else if (command === "forbidden7Days") {
    await forbidden(7);
  } else if (command === "forbiddenForever") {
    await forbidden(-1);
  } else {
    console.log("click on item " + command);
  }
}
async function forbidden(days) {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: props.article.user.id,
        days,
      },
    });
    useMsgSuccess(i18n.t('alert.mute_user_success'));
  } catch (e) {
    useMsgError(i18n.t('alert.mute_user_failure'));
  }
}
function deleteArticle() {
  useConfirm(i18n.t('dialog.message.confirm_delete_article')).then(function () {
    useHttpPost(`/api/article/delete/${props.article.id}`)
      .then(() => {
        useMsg({
          message: i18n.t('alert.delete_success'),
          onClose() {
            useLinkTo("/articles");
          },
        });
      })
      .catch((e) => {
        useMsgError(i18n.t('alert.delete_success', { error: (e.message || e) }));
      });
  });
}
function editArticle() {
  useLinkTo(`/article/edit/${props.article.id}`);
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--text-color3);
  font-size: 12px;
}

.action-menu {
  ul {
    li {
      text-transform: capitalize;
    }
  }
}
</style>
