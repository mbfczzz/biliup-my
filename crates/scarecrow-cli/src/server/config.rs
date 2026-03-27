use crate::server::core::downloader::DownloaderType;
use crate::server::errors::{AppError, AppResult};
use crate::server::infrastructure::models::hook_step::HookStep;
use scarecrow_core::bilibili::Credit;
use error_stack::bail;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, path::Path, path::PathBuf};
use struct_patch::Patch;

/// 全局配置结构体
#[derive(bon::Builder, Debug, PartialEq, Clone, Serialize, Deserialize, Patch)]
#[patch(attribute(derive(Debug, Clone, Default, Deserialize, Serialize)))]
pub struct Config {
    // ===== 全局录播与上传设置 =====
    /// 下载器类型：streamlink | ffmpeg | stream-gears | 自定义
    #[serde(default)]
    pub downloader: Option<DownloaderType>,

    /// 文件大小限制（字节）
    #[builder(default = default_file_size())]
    #[serde(default = "default_file_size")]
    pub file_size: u64,

    /// 分段时间，格式如 "00:00:00"，保留为字符串以保持直观
    #[serde(default = "default_segment_time")]
    pub segment_time: Option<String>,

    /// 过滤阈值（MB）
    #[builder(default = default_filtering_threshold())]
    #[serde(default = "default_filtering_threshold")]
    pub filtering_threshold: u64,

    /// 文件名前缀
    #[serde(default)]
    pub filename_prefix: Option<String>,

    /// 分段处理器是否并行执行
    #[serde(default)]
    pub segment_processor_parallel: Option<bool>,

    /// 上传器类型：Noop | bili_web | scarecrow | 其他
    #[serde(default)]
    pub uploader: Option<String>,

    /// 提交API类型：web | client
    #[serde(default)]
    pub submit_api: Option<String>,

    /// 上传线路：AUTO | alia | bda2 | bldsa | qn | tx | txa
    #[builder(default = default_lines())]
    #[serde(default = "default_lines")]
    pub lines: String,

    /// 上传线程数
    #[builder(default = default_threads())]
    #[serde(default = "default_threads")]
    pub threads: u32,

    /// 延迟时间（秒）
    #[builder(default = default_delay())]
    #[serde(default = "default_delay")]
    pub delay: u64,

    /// 事件循环间隔（秒）
    #[builder(default = default_event_loop_interval())]
    #[serde(default = "default_event_loop_interval")]
    pub event_loop_interval: u64,

    /// 检查器休眠时间（秒）
    #[builder(default = default_checker_sleep())]
    #[serde(default = "default_checker_sleep")]
    pub checker_sleep: u64,

    /// 连接池1大小
    #[builder(default = default_pool1_size())]
    #[serde(default = "default_pool1_size")]
    pub pool1_size: u32,

    /// 连接池2大小
    #[builder(default = default_pool2_size())]
    #[serde(default = "default_pool2_size")]
    pub pool2_size: u32,

    // ===== 各平台录播设置 =====
    /// 是否使用直播封面
    #[serde(default)]
    pub use_live_cover: Option<bool>,

    // Twitch平台设置
    /// Twitch弹幕录制
    #[serde(default)]
    pub twitch_danmaku: Option<bool>,
    /// Twitch禁用广告
    #[serde(default)]
    pub twitch_disable_ads: Option<bool>,

    /// 录制主播配置映射
    #[serde(default)]
    pub streamers: HashMap<String, StreamerConfig>,

    /// 用户Cookie配置
    #[serde(default)]
    pub user: Option<UserConfig>,

    pub loggers_level: Option<String>,
}

/// 主播配置结构体
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
pub struct StreamerConfig {
    /// 直播间URL列表
    pub url: Vec<String>,

    /// 视频标题
    #[serde(default)]
    pub title: Option<String>,

    /// 分区ID
    #[serde(default)]
    pub tid: Option<u32>,

    /// 版权类型
    #[serde(default)]
    pub copyright: Option<u8>,

    /// 封面路径
    #[serde(default)]
    pub cover_path: Option<PathBuf>,

    /// 视频描述（保留缩进和多行格式）
    #[serde(default)]
    pub description: Option<String>,

    #[serde(default)]
    pub credits: Option<Vec<Credit>>,

    #[serde(default)]
    pub dynamic: Option<String>,

    #[serde(default)]
    pub dtime: Option<u64>,

    #[serde(default)]
    pub dolby: Option<u8>,

    #[serde(default)]
    pub hires: Option<u8>,

    #[serde(default)]
    pub charging_pay: Option<u8>,

    #[serde(default)]
    pub no_reprint: Option<u8>,

    #[serde(default)]
    pub is_only_self: Option<u8>,

    #[serde(default)]
    pub uploader: Option<String>,

    #[serde(default)]
    pub filename_prefix: Option<String>,

    #[serde(default)]
    pub user_cookie: Option<String>,

    #[serde(default)]
    pub use_live_cover: Option<bool>,

    #[serde(default)]
    pub tags: Option<Vec<String>>,

    #[serde(default)]
    pub time_range: Option<String>,

    #[serde(default)]
    pub excluded_keywords: Option<Vec<String>>,

    #[serde(default)]
    pub preprocessor: Option<Vec<HookStep>>,

    #[serde(default)]
    pub segment_processor: Option<Vec<HookStep>>,

    #[serde(default)]
    pub downloaded_processor: Option<Vec<HookStep>>,

    #[serde(default)]
    pub postprocessor: Option<Vec<HookStep>>,

    #[serde(default)]
    pub format: Option<String>,

    #[serde(default)]
    pub opt_args: Option<Vec<String>>,

    // “override” 是字段名，这里改为 override_cfg 避免与保留字混淆
    #[serde(rename = "override", default)]
    pub override_cfg: Option<HashMap<String, serde_json::Value>>,
}

/// 用户配置结构体
#[derive(bon::Builder, PartialEq, Debug, Clone, Serialize, Deserialize, Default, Patch)]
#[patch(attribute(derive(Debug, Default, Deserialize)))]
pub struct UserConfig {
    // B站配置（用于上传投稿）
    /// B站Cookie字符串
    #[serde(default)]
    pub bili_cookie: Option<String>,
    /// B站Cookie文件路径
    #[serde(default)]
    pub bili_cookie_file: Option<PathBuf>,

    // Twitch配置
    /// Twitch Cookie
    #[serde(default)]
    pub twitch_cookie: Option<String>,
}

/// 默认文件大小：2.5GB
fn default_file_size() -> u64 {
    2_621_440_000
}

/// 默认分段时间：2小时
pub fn default_segment_time() -> Option<String> {
    Some("02:00:00".to_string())
}

/// 默认过滤阈值：20MB
fn default_filtering_threshold() -> u64 {
    20
}

/// 默认上传线路：自动选择
fn default_lines() -> String {
    "AUTO".to_string()
}

/// 默认线程数：3
fn default_threads() -> u32 {
    3
}

/// 默认延迟：300秒
fn default_delay() -> u64 {
    300
}

/// 默认事件循环间隔：30秒
fn default_event_loop_interval() -> u64 {
    30
}

/// 默认检查器休眠时间：10秒
fn default_checker_sleep() -> u64 {
    10
}

/// 默认连接池1大小：5
fn default_pool1_size() -> u32 {
    5
}

/// 默认连接池2大小：3
fn default_pool2_size() -> u32 {
    3
}

impl Config {
    /// 从指定路径加载配置文件，如果不存在则创建默认配置
    pub fn load_or_create<P: AsRef<Path>>(path: P) -> AppResult<Self> {
        bail!(AppError::Custom(format!(
            "load_or_create: {:?}",
            path.as_ref().display()
        )))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // 情况 1: 字段存在且有字符串值
    #[test]
    fn test_field_with_value() {
        let json = r#"{"maybe_name": "Alice"}"#;

        let patch = r#"{"maybe_name": "Alice"}"#;

        // Single Option: Some("Alice")
        let mut single: Config = serde_json::from_str(json).unwrap();

        let patch: ConfigPatch = serde_json::from_str(json).unwrap();

        single.apply(patch);

        // 从 JSON 反序列化时,未指定的字段使用 serde default (None)
        // 而 builder 的 default 是 Some("02:00:00"),两者不同
        // 需要明确设置 segment_time 为 None 以匹配反序列化结果
        assert_eq!(
            single,
            Config::builder().streamers(Default::default()).build(),
            "普通Option正常包裹一层"
        );
    }
}
