# API 文档

本文档描述了基于 Gin 框架构建的 API 接口，包括登录、查询、签到和学校搜索等功能。所有接口均位于 `/api` 路径下。

## 目录

- [基础信息](#基础信息)
- [API 列表](#api-列表)
    - [1. 用户登录](#1-用户登录)
    - [2. 查询用户信息](#2-查询用户信息)
    - [3. 用户签到](#3-用户签到)
    - [4. 搜索学校](#4-搜索学校)
- [错误码](#错误码)
- [示例](#示例)

---

## 基础信息

- **基座 URL**: `http://<your-domain>/api`
- **数据格式**: JSON
- **认证方式**: 部分接口可能需要 Token 认证（具体视实现而定）

---

## API 列表

### 1. 用户登录

**Endpoint**

```
POST /api/login
```

**描述**

用户通过账号、密码和学校 ID 进行登录，成功后返回一个认证 Token。

**请求头**

```http
Content-Type: application/json
```

**请求体**

| 字段      | 类型   | 描述                         | 必填 |
| --------- | ------ |----------------------------| ---- |
| account   | string | 用户账号                       | 是   |
| password  | string | 用户密码                       | 是   |
| school_id | string | 学校唯一标识（可使用search_school查询） | 是   |
| token     | string | （可选）额外 Token               | 否   |

**示例请求**

```json
{
  "account": "user123",
  "password": "securepassword",
  "school_id": "school_001",
  "token": "optional_token"
}
```

**响应**

- **成功**

  ```json
  {
    "message": "登录成功",
    "token": "generated_auth_token"
  }
  ```

- **失败**

  ```json
  {
    "message": "登录失败",
    "error": "错误描述"
  }
  ```

**状态码**

- `200 OK`：登录成功
- `400 Bad Request`：请求参数错误
- `401 Unauthorized`：登录失败

---

### 2. 查询用户信息

**Endpoint**

```
GET /api/query
```

**描述**

根据用户账号查询相关信息。

**请求参数**

| 参数    | 类型   | 描述       | 必填 |
| ------- | ------ | ---------- | ---- |
| account | string | 用户账号   | 是   |

**示例请求**

```
GET /api/query?account=user123
```

**响应**

- **成功**

  ```json
    {
    "message": "登录成功",
    "token": "6e5dc883fe6855b8fd6616900624c9a3"
    }
  ```

- **失败**

一共七次登录密码尝试的机会

  ```json
    {
    "message": "登录失败",
    "error": "登录失败: 登录失败,还剩6次机会"
    }
  ```

**状态码**

- `200 OK`：查询成功
- `400 Bad Request`：缺少必要参数
- `500 Internal Server Error`：查询失败

---

### 3. 用户签到

**Endpoint**

```
POST /api/sign
```

**描述**

用户进行签到操作，记录签到信息。

**请求头**

```http
Content-Type: application/json
```

**请求体**

| 字段        | 类型   | 描述            | 必填 |
| ----------- | ------ | --------------- | ---- |
| account     | string | 用户账号        | 是   |
| address     | string | 签到地址        | 是   |
| address_name | string | 地址名称        | 否   |
| latitude    | string | 纬度            | 否   |
| longitude   | string | 经度            | 否   |
| remark      | string | 备注            | 否   |
| comment     | string | 评论            | 否   |
| province    | string | 省份            | 否   |
| city        | string | 城市            | 否   |
| secret_key  | string | （可选）密钥    | 否   |

**示例请求**

```json
{
  "account": "user123",
  "address": "北京市朝阳区",
  "address_name": "公司",
  "latitude": "39.9042",
  "longitude": "116.4074",
  "remark": "0",
  "comment": "我在上班",
  "province": "北京市",
  "city": "朝阳区",
  "secret_key": "" 
}
```

**响应**

- **成功**

  ```json
  {
    "message": "签到成功"
  }
  ```

- **失败**

  ```json
  {
    "message": "签到失败",
    "error": "错误描述"
  }
  ```

**状态码**

- `200 OK`：签到成功
- `400 Bad Request`：请求参数错误
- `500 Internal Server Error`：签到失败

---

### 4. 搜索学校

**Endpoint**

```
GET /api/search_school
```

**描述**

根据学校名称搜索对应的学校信息。

**请求参数**

| 参数        | 类型   | 描述       | 必填 |
| ----------- | ------ | ---------- | ---- |
| school_name | string | 学校名称   | 是   |

**示例请求**

```
GET /api/search_school?school_name=南京
```

**响应**

- **成功**

  ```json
    {
  "message": "查询成功",
  "schools": [
    {
      "school_id": "1690",
      "school_name": "南京中医药大学翰林学院"
    },
    {
      "school_id": "529",
      "school_name": "南京中华中等专业学校"
    },
    {
      "school_id": "714",
      "school_name": "南京交通技师学院"
    },
    {
      "school_id": "1777",
      "school_name": "南京交通职业技术学院"
    },
    {
      "school_id": "1745",
      "school_name": "南京信息工程大学"
    },
    {
      "school_id": "2218",
      "school_name": "南京信息工程技术学校"
    },
    {
      "school_id": "1580",
      "school_name": "南京信息职业技术学院"
    },
    {
      "school_id": "467",
      "school_name": "南京六合中等专业学校"
    },
    {
      "school_id": "1787",
      "school_name": "南京化工职业技术学院"
    },
    {
      "school_id": "1746",
      "school_name": "南京医科大学康达学院"
    },
    {
      "school_id": "1748",
      "school_name": "南京工业大学"
    },
    {
      "school_id": "1653",
      "school_name": "南京工业大学浦江学院"
    },
    {
      "school_id": "1756",
      "school_name": "南京工业职业技术大学"
    },
    {
      "school_id": "1557",
      "school_name": "南京市莫愁中等专业学校"
    },
    {
      "school_id": "2257",
      "school_name": "南京师范大学"
    },
    {
      "school_id": "522",
      "school_name": "南京技师学院"
    },
    {
      "school_id": "193",
      "school_name": "南京旅游职业学院"
    },
    {
      "school_id": "2258",
      "school_name": "南京晓庄学院"
    },
    {
      "school_id": "1785",
      "school_name": "南京机电职业技术学院"
    },
    {
      "school_id": "1769",
      "school_name": "南京森林公安高等专科学校"
    },
    {
      "school_id": "720",
      "school_name": "南京江宁高等职业技术学校"
    },
    {
      "school_id": "1767",
      "school_name": "南京特殊教育职业技术学院"
    },
    {
      "school_id": "354",
      "school_name": "南京航空航天大学金城学院"
    },
    {
      "school_id": "508",
      "school_name": "南京金陵中等专业学校"
    },
    {
      "school_id": "1791",
      "school_name": "南京铁道职业技术学院"
    },
    {
      "school_id": "1030",
      "school_name": "南京高等职业技术学校"
    },
    {
      "school_id": "1741",
      "school_name": "江苏省南京工程高等职业学校"
    },
    {
      "school_id": "1795",
      "school_name": "民办明达职业技术学院南京市"
    }
  ]
    }
  ```

- **未找到**

  ```json
  {
    "message": "没有找到匹配的学校"
  }
  ```

- **失败**

  ```json
  {
    "message": "查询失败",
    "error": "错误描述"
  }
  ```

**状态码**

- `200 OK`：查询成功
- `400 Bad Request`：缺少必要参数
- `404 Not Found`：没有找到匹配的学校
- `500 Internal Server Error`：查询失败

---

## 错误码

| 状态码 | 描述             |
| ------ | ---------------- |
| 400    | 请求参数错误     |
| 401    | 未授权（登录失败）|
| 404    | 资源未找到       |
| 500    | 服务器内部错误   |

---

## 示例

### 登录示例

**请求**

```http
POST /api/login
Content-Type: application/json

{
  "account": "user123",
  "password": "securepassword",
  "school_id": "0"
}
```

**响应**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "登录成功",
  "token": "abcdef123456"
}
```

### 查询用户信息示例

**请求**

```
GET /api/query?account=user123
```

**响应**

```json
    {
  "message": "查询成功",
  "data": {
    "banner_list": [
      {
        "pic_url": "https://example.com/banner1.jpg",
        "redirect_url": ""
      },
      {
        "pic_url": "https://example.com/banner2.jpg",
        "redirect_url": ""
      }
    ],
    "continuous_sign_in": "48",
    "current_company": "",
    "first_time_sign": 0,
    "mark_list": [
      { "key": 0, "value": "上班" },
      { "key": 1, "value": "因公外出" },
      { "key": 2, "value": "假期" },
      { "key": 3, "value": "请假" },
      { "key": 4, "value": "轮岗" },
      { "key": 5, "value": "回校" },
      { "key": 6, "value": "外宿" },
      { "key": 7, "value": "在家" },
      { "key": 8, "value": "下班" },
      { "key": 9, "value": "学习" },
      { "key": 10, "value": "毕业设计" },
      { "key": 11, "value": "院区轮转" },
      { "key": 13, "value": "集训" },
      { "key": 14, "value": "休息" }
    ],
    "offset_circle_building": 1500,
    "offset_distance": 1000,
    "rule": [
      "签到规则：",
      "· 每日签到表现积分增加2分，异常签到1分。",
      "· 连续签到5天奖励表现积分5分。",
      "· 签到地址和实习地址需在首次签到时进行确认。"
    ],
    "sign_in_month": [
      { "sign_time": "1735596959", "sign_time_text": "2024-12-31", "status_code": "0" },
      { "sign_time": "1735510546", "sign_time_text": "2024-12-30", "status_code": "0" },
      { "sign_time": "1735424146", "sign_time_text": "2024-12-29", "status_code": "0" }
    ],
    "sign_resources_info": {
      "mid_sign_address": "xx省某市某区某一路XX号靠近某工业园",
      "mid_sign_latitude": "xx.xxxxxx",
      "mid_sign_longitude": "xx.xxxxxx",
      "mid_sign_time": "2024-11-30 08:21:41"
    }
  }
}

```

---

如有任何疑问或建议，请联系开发团队以获取支持。