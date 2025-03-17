<template>
  <div class="left-container">
    <my-counts :user="localUser" />
    <my-profile :user="localUser" />
    <fans-widget :user="localUser" />
    <follow-widget :user="localUser" />

    <div v-if="isAdmin" class="widget">
      <div class="widget-header">{{ $t('profile.actions.title') }}</div>
      <div class="widget-content">
        <ul class="operations">
          <li v-if="localUser.forbidden">
            <icon name="Ban" />
            <a @click="removeForbidden">&nbsp;{{ $t('profile.actions.unmute') }}</a>
          </li>
          <template v-else>
            <li>
              <icon name="Ban" />
              <a @click="forbidden(7)">&nbsp;{{ $t("profile.actions.mute_7days") }}</a>
            </li>
            <li v-if="isSiteOwner">
              <icon name="Ban" />
              <a @click="forbidden(-1)">&nbsp;{{ $t("profile.actions.mute_permanent") }}</a>
            </li>
          </template>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n();
import { ElMessageBox } from "element-plus";
const userStore = useUserStore();
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});
const localUser = ref(props.user);

const isSiteOwner = computed(() => {
  return userIsOwner(userStore.user);
});

const isAdmin = computed(() => {
  return userIsOwner(userStore.user) || userIsAdmin(userStore.user);
});

function forbidden(days) {
  const msg = days > 0 ? i18n.t('dialog.message.mute_user') : i18n.t('dialog.message.mute_user_permanent');
  ElMessageBox.confirm(msg, i18n.t('dialog.title.prompt'), {
    confirmButtonText: i18n.t('dialog.button.confirm'),
    cancelButtonText: i18n.t('dialog.button.cancel'),
    type: "warning",
  })
    .then(() => {
      doForbidden(days);
    })
    .catch(() => { });
}

async function doForbidden(days) {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: localUser.value.id,
        days,
      },
    });
    localUser.value.forbidden = true;
    useMsgSuccess(i18n.t('alert.mute_user_success'));
  } catch (e) {
    useMsgError(i18n.t('alert.mute_user_failure'));
  }
}

async function removeForbidden() {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: localUser.value.id,
        days: 0,
      },
    });
    localUser.value.forbidden = false;
    useMsgSuccess(i18n.t('alert.unumute_user_success'));
  } catch (e) {
    useMsgError(i18n.t('alert.unumute_user_failure'));
  }
}
</script>

<style lang="scss" scoped>
.img-avatar {
  margin-top: 5px;
  border: 1px dotted var(--border-color);
  border-radius: 5%;
}

.operations {
  list-style: none;

  li {
    padding-left: 3px;

    font-size: 13px;

    &:hover {
      cursor: pointer;
      background-color: #fcf8e3;
      color: #8a6d3b;
      font-weight: bold;
    }

    a {
      color: var(--text-link-color);
    }
  }
}
</style>
