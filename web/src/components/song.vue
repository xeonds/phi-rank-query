<script setup lang="ts">
const props = defineProps({
    index: {
        type: String,
        default: '0'
    },
    song: {
        type: Object,
        default: {
        }
    }
});
const index = props.index ?? '#?';
const data = props.song ?? {
    illustration: '',
    rank: '??',
    difficulty: '??',
    rks: '??',
    song: '??',
    Rating: '??',
    score: '??',
    acc: '??',
}
</script>

<template>
    <div class="song">
        <div class="ill-box">
            <div class="num">
                <p name="pvis">{{ index }}</p>
            </div>
            <div class="ill">
                <img v-lazy="data.illustration" alt="ill">
            </div>
            <div :class="`rank-${data.rank}`">
                <div class="org">
                    <p>{{ data.rank }}&ensp;{{ parseFloat(data.difficulty).toFixed(1) }}</p>
                </div>
                <div class="rel">
                    <p>{{ parseFloat(data.rks).toFixed(2) }}</p>
                </div>
            </div>
        </div>
        <div class="info">
            <div class="songname">
                <p name="pvis">{{ data.song }}</p>
            </div>
            <div class="songinfo">
                <div class="Rating">
                    <img :src="'/assets/' + data.Rating + '.png'" alt="Rating">
                </div>
                <div class="chengji">
                    <div class="score">
                        <p>{{ data.score }}</p>
                    </div>
                    <div class="line"></div>
                    <div class="acc-box">
                        <div class="acc">
                            <p>{{ parseFloat(data.acc).toFixed(2) }}%</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.song {
    width: auto;
    height: 10.5rem;
    display: flex;
    flex-flow: row;
    align-items: center;
    justify-content: center;

    .ill-box {
        width: auto;
        height: inherit;
        display: flex;
        flex-flow: column;
        align-items: flex-start;
        justify-content: center;
        z-index: 1;
        filter: drop-shadow(2px 4px 10px #00000055);

        .num {
            p {
                font-size: 0.8rem;
                margin: 0;
                color: black;
            }

            z-index: 2;
            width: 1.8rem;
            height: 1.2rem;
            display: flex;
            background-color: rgba(255, 255, 255, 1);
            clip-path: polygon(100% 0, 90% 100%, 0 100%, 10% 0);
            justify-content: center;
            align-items: center;
        }

        .ill {
            width: 10rem;
            height: auto;
            display: flex;
            flex-direction: column;
            margin-top: -1.2rem;
            margin-bottom: -2.5rem;
            clip-path: polygon(100% 0, 90% 100%, 0 100%, 10% 0);
            overflow: hidden;
        }

        .rank-EZ {
            background-color: #92D050;
        }

        .rank-HD {
            background-color: #00B0F0;
        }

        .rank-IN {
            background-color: rgb(255, 0, 0);
        }

        .rank-AT {
            background-color: rgb(110, 110, 110);
        }

        .rank-EZ,
        .rank-HD,
        .rank-IN,
        .rank-AT {
            width: 4rem;
            height: 2.5rem;
            display: flex;
            align-items: center;
            flex-flow: column;
            justify-content: center;
            margin-left: -15%;
            clip-path: polygon(100% 0, 90% 100%, 0 100%, 10% 0);

            .org,
            .rel {
                display: flex;
                width: auto;
                justify-content: center;
                align-items: center;
            }

            .org {
                height: 35%;
                margin-left: 10%;
                font-size: 0.8rem;
            }

            .rel {
                height: 49%;
                margin-right: 10%;
                font-size: 1rem;
            }
        }
    }

    .info {
        margin-left: -2rem;
        padding-left: 2rem;
        padding-top: .2rem;
        padding-bottom: .2rem;
        display: flex;
        flex-flow: column;
        align-items: center;
        justify-content: center;
        width: 14rem;
        height: 5rem;
        background-color: rgba(0, 0, 0, 0.5);
        clip-path: polygon(100% 0, 95% 100%, 0 100%, 5% 0);

        .songname {
            width: 90%;
            height: 25%;
            display: flex;
            justify-content: center;
            align-items: center;

            p {
                font-size: .75rem;
                margin: 0;
                text-align: center;
            }
        }

        .songinfo {
            width: 100%;
            height: 70%;
            display: flex;
            flex-flow: row;
            align-items: center;

            .Rating {
                width: 2.5rem;
                height: 2.5rem;
                display: flex;

                img {
                    width: 100%;
                    filter: drop-shadow(0px 0px 10px #ffffffbb);
                }
            }

            .chengji {
                width: 80%;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;

                .score {
                    font-size: 1.4rem;
                }

                .line {
                    width: 90%;
                    margin: 0;
                    border-bottom: 1px solid;
                    border-color: white;
                }

                .acc-box {
                    display: flex;
                    justify-content: center;
                    align-items: flex-start;

                    .acc {
                        display: flex;
                        justify-content: flex-end;
                        align-items: center;

                        p {
                            font-size: .9rem;
                        }
                    }

                    .suggest {
                        display: none;
                        width: auto;
                        display: flex;

                        p {
                            font-size: 1rem;
                            color: rgb(192, 221, 173);
                        }
                    }
                }
            }
        }
    }
}
</style>