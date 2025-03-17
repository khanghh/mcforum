<script setup>
// const configStore = useConfigStore()
useHead({
  title: 'About',
})
</script>

<template>
  <section class="main">
    <div class="container">
      <div class="widget">
        <div class="widget-header">
          About
        </div>
        <div class="widget-content content">
          <h2>Introduction</h2>
          <p>
            BBS-GO is a beautifully crafted, lightweight, open-source community forum system developed in Go language.
            Initially, this project was used for learning and communication. After being open-sourced, more and more
            people began to appreciate and follow it, which has become the motivation for me to continue upgrading and
            maintaining it.
          </p>

          <p>
            To facilitate learning and communication, we used BBS-GO to set up <a
              href="https://mlog.club"
              target="_blank">MLog Club</a> in August 2019. Everyone is welcome to join the discussion.
          </p>

          <h2>Contact Me</h2>
          <p>
            If you have any thoughts, suggestions, or questions, feel free to visit <a
              href="https://mlog.club"
              target="_blank">MLog Club</a> for communication.
          </p>
          <h2>Project Contributors</h2>
          <p>
            <a
              href="https://github.com/mlogclub/bbs-go/graphs/contributors"
              target="_blank">Contributors</a>
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped></style>
