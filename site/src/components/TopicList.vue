<template>
  <ul class="topic-list">
    <li v-for="topic in topics" :key="topic.id" class="topic-item">
      <div class="topic-avatar" :title="topic.user.nickname">
        <my-avatar :user="topic.user" />
      </div>
      <div class="topic-main-content">
        <div class="topic-userinfo">
          <div class="infos">
            <my-avatar
              class="topic-inline-avatar"
              :user="topic.user"
              :size="20" />
            <nuxt-link :to="`/user/${topic.user.id}`" class="topic-nickname">
              {{ topic.user.nickname }}
            </nuxt-link>
          </div>
          <div class="icons">
            <span v-if="showPinned && topic.pinned" class="topic-pinned-icon">{{ $t('feed.pinned') }}</span>
          </div>
        </div>
        <div class="topic-time">
          {{ $t('feed.published_on') }} {{ usePrettyDate(topic.createTime) }}
        </div>
        <div class="topic-content" :class="{ 'topic-tweet': topic.type === 1 }">
          <template v-if="topic.type === 0">
            <h1 class="topic-title">
              <nuxt-link :to="`/topic/${topic.id}`">
                {{ topic.title }}
              </nuxt-link>
            </h1>
            <nuxt-link
              :to="`/topic/${topic.id}`"
              class="topic-summary">
              {{ topic.summary }}
            </nuxt-link>
          </template>
          <template v-if="topic.type === 1">
            <nuxt-link
              v-if="topic.content"
              :to="`/topic/${topic.id}`"
              class="topic-summary">
              {{ topic.content }}
            </nuxt-link>
            <ul
              v-if="topic.imageList && topic.imageList.length"
              class="topic-image-list">
              <li v-for="(image, index) in topic.imageList" :key="index">
                <nuxt-link
                  :to="`/topic/${topic.id}`"
                  class="image-item">
                  <img :src="image.preview" />
                </nuxt-link>
              </li>
            </ul>
          </template>
        </div>
        <div class="topic-bottom">
          <div class="topic-handlers">
            <div class="btn" :class="{ liked: topic.liked }" @click="like(topic)">
              <icon name="ThumbsUp" :filled="topic.liked" />
              {{ topic.liked ? $t('feed.actions.liked') : $t('feed.actions.like') }}
              <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.id)">
              <icon name="MessageSquareMore" size="1em" />
              {{ $t('feed.actions.comment') }}
              <span v-if="topic.commentCount > 0">
                {{ topic.commentCount }}
              </span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.id)">
              <icon name="BookOpenText" size="1em" />
              {{ $t('feed.actions.view') }}
              <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
            </div>
          </div>
          <div class="topic-tags">
            <nuxt-link v-if="topic.node" class="topic-tag" :to="`/f/${topic.node.slug}`" :alt="topic.node.name">
              {{ topic.node.name }}
            </nuxt-link>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script setup>
const i18n = useI18n()
defineProps({
  topics: {
    type: Array,
    default: () => [],
    required: false,
  },
  showAvatar: {
    type: Boolean,
    default: true,
  },
  showPinned: {
    type: Boolean,
    default: false,
  },
})

async function like(topic) {
  try {
    if (topic.liked) {
      await useHttpPostForm('/api/like/unlike', {
        body: {
          entityType: 'topic',
          entityId: topic.id,
        },
      })
      topic.liked = false
      topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    }
    else {
      await useHttpPostForm('/api/like/like', {
        body: {
          entityType: 'topic',
          entityId: topic.id,
        },
      })
      topic.liked = true
      topic.likeCount++
      useMsgSuccess(i18n.t('message.liked_success'))
    }
  }
  catch (e) {
    useCatchError(e)
  }
}

async function toTopicDetail(topicId) {
  useLinkTo(`/topic/${topicId}`)
}
</script>

<style lang="scss" scoped>
.topic-list {
  .topic-item {
    padding: 12px 12px;
    display: flex;
    position: relative;
    overflow: hidden;
    transition: background-color 0.5s;
    border-radius: 3px;
    background: var(--bg-color);

    &:not(:last-child) {
      margin-bottom: 10px;
    }

    &:hover {
      background: var(--bg-color2);
    }

    .topic-avatar {
      //
    }

    .topic-main-content {
      flex: 1;
      margin-left: 12px;

      // .topic-top {
      //   margin-bottom: 8px;

      //   .topic-userinfo {
      //     display: inline-flex;
      //     align-items: center;

      //     .topic-nickname {
      //       // font-weight: 700;
      //       font-size: 14px;
      //       color: var(--text-color);
      //       display: flex;
      //       max-width: 250px;
      //       overflow: hidden;
      //     }

      //     .topic-inline-avatar {
      //       display: none;
      //       margin-right: 5px;
      //     }

      //     .topic-sticky-icon {
      //       color: var(--color-red);
      //       border: 1px solid var(--color-red);
      //       border-radius: 2px;
      //       font-size: 12px;
      //       font-weight: 700;
      //       padding: 0 5px;
      //       margin-left: 10px;
      //     }
      //   }

      //   .topic-time {
      //     color: var(--text-color3);
      //     font-size: 12px;
      //     float: right;
      //     display: flex;
      //   }

      //   @media screen and (max-width: 1024px) {
      //     .topic-time {
      //       float: none;
      //       margin-top: 8px;
      //     }
      //   }
      // }

      .topic-userinfo {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .infos {
          flex: 1;
          display: flex;

          .topic-nickname {
            font-size: 15px;
            color: var(--text-color);
            display: flex;
            max-width: 250px;
            overflow: hidden;
          }

          .topic-inline-avatar {
            display: none;
            margin-right: 5px;
          }
        }

        .icons {
          display: flex;

          .topic-pinned-icon {
            color: var(--color-red);
            border: 1px solid var(--color-red);
            border-radius: 2px;
            font-size: 12px;
            font-weight: 500;
            padding: 0 5px;
            margin-left: 10px;
          }
        }
      }

      .topic-time {
        margin-top: 2px;
        color: var(--text-color3);
        font-size: 12px;
      }

      .topic-content {
        margin-top: 6px;

        .topic-title {
          display: inline-block;
          margin-bottom: 6px;
          word-wrap: break-word;
          word-break: break-all;
          width: 100%;

          a {
            font-size: 16px;
            font-weight: 400;
            color: var(--text-color);

            &:hover {
              //color: #3273dc;
              text-decoration: underline;
            }
          }
        }

        .topic-summary {
          font-size: 14px;
          margin-bottom: 6px;
          width: 100%;
          text-decoration: none;
          color: var(--text-color3);
          word-wrap: break-word;

          overflow: hidden;
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 3;
          text-align: justify;
          word-break: break-all;
          text-overflow: ellipsis;
        }

        &.topic-tweet {
          .topic-summary {
            color: var(--text-color);
            white-space: pre-line;
          }
        }

        .topic-image-list {
          margin-top: 10px;

          li {
            cursor: pointer;
            text-align: center;

            display: inline-block;
            vertical-align: middle;
            margin: 0 8px 8px 0;
            background-color: var(--bg-color2);
            background-size: 32px 32px;
            background-position: 50%;
            background-repeat: no-repeat;
            overflow: hidden;
            position: relative;

            .image-item {
              display: block;
              overflow: hidden;
              transform-style: preserve-3d;

              &>img {
                width: 100%;
                height: 100%;
                object-fit: cover;
                transition: all 0.5s ease-out 0.1s;

                &:hover {
                  transform: matrix(1.04, 0, 0, 1.04, 0, 0);
                  backface-visibility: hidden;
                }
              }
            }

            /* 只有一个图片时 */
            &:first-child:nth-last-child(1) {
              width: 210px;
              height: 210px;
              line-height: 210px;

              .image-item {
                width: 210px;
                height: 210px;
              }
            }

            /* 只有两个图片时 */
            &:first-child:nth-last-child(2),
            &:first-child:nth-last-child(2)~li {
              width: 180px;
              height: 180px;
              line-height: 180px;

              .image-item {
                width: 180px;
                height: 180px;
              }
            }

            /*大于两个图片时*/
            &:first-child:nth-last-child(n + 3),
            &:first-child:nth-last-child(n + 3)~li {
              width: 120px;
              height: 120px;
              line-height: 120px;

              .image-item {
                width: 120px;
                height: 120px;
              }
            }
          }
        }
      }

      .topic-bottom {
        display: flex;

        .topic-handlers {
          display: flex;
          align-items: center;
          margin-top: 6px;
          font-size: 12px;
          flex: 1;
          user-select: none;

          .btn {
            color: var(--text-color3);
            cursor: pointer;

            &:not(:last-child) {
              margin-right: 20px;
            }

            &:hover {
              color: var(--text-link-color);
            }

            i {
              margin-right: 3px;
              font-size: 12px;
              position: relative;
            }
          }
        }

        .topic-tags {
          .topic-tag {
            padding: 2px 8px;
            justify-content: center;
            align-items: center;
            border-radius: 12.5px;
            margin-right: 10px;
            background: var(--bg-color2);
            border: 1px solid var(--border-color2);
            color: var(--text-color3);
            font-size: 12px;

            &:hover {
              color: var(--text-link-color);
              background: var(--bg-color);
            }
          }
        }
      }

      .liked {
        color: var(--color-blue) !important;
      }
    }

    @media screen and (max-width: 768px) {
      .topic-avatar {
        display: none;
      }

      .topic-main-content {
        margin-left: 0;

        .topic-inline-avatar {
          display: block !important;
        }

        .topic-bottom {
          .topic-tags {
            display: none;
          }
        }
      }
    }
  }
}
</style>
