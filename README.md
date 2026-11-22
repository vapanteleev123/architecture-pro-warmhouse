# Project_template

Это шаблон для решения проектной работы. Структура этого файла повторяет структуру заданий. Заполняйте его по мере работы над решением.

# Задание 1. Анализ и планирование

<aside>

Чтобы составить документ с описанием текущей архитектуры приложения, можно часть информации взять из описания компании и условия задания. Это нормально.

</aside>

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи могут удалённо включать и выключать отопление в своих домах через веб-интерфейс.
- Система управляет модулями отопления (датчики и реле): команды идут синхронно от сервера к устройству.
- Подключение нового модуля выполняется специалистом компании; пользователь самостоятельно подключить датчик не может.
- При включении/подключении датчики регистрируются в системе «Умный дом».

**Мониторинг температуры:**

- Система получает данные о температуре с датчиков, установленных в домах.
- Опрос температуры выполняется синхронно: сервер запрашивает показания у датчика/внешнего temperature API.
- Пользователи могут просматривать текущую температуру в своих домах через веб-интерфейс.
- В монолите поддерживается CRUD по датчикам (создание, чтение, обновление, удаление) и обновление текущего значения/статуса датчика.

### 2. Анализ архитектуры монолитного приложения

| Аспект | As-Is |
| --- | --- |
| Язык программирования | Go (Gin) |
| База данных | PostgreSQL |
| Архитектура | Монолит: HTTP-обработка, бизнес-логика и доступ к данным в одном приложении (`smart_home`) |
| Взаимодействие | Синхронное (запрос–ответ); асинхронных вызовов и брокеров сообщений нет |
| Внешние зависимости | Temperature API (синхронный HTTP) для актуальных показаний температуры |
| Масштабируемость | Ограничена: нельзя независимо масштабировать отдельные функции (управление устройствами, телеметрия, отопление) |
| Развёртывание | Требует остановки/перезапуска всего приложения при любом изменении |
| Клиенты / нагрузка (сейчас) | ~100 веб-клиентов, ~100 модулей управления отоплением |

Компоненты монолита внутри одного процесса: маршрутизация API (`/api/v1/sensors`), handlers, сервис температуры, слой работы с PostgreSQL.

### 3. Определение доменов и границы контекстов

На основе функциональности As-Is и целевой экосистемы выделены следующие **ограниченные контексты (Bounded Contexts)** в духе DDD:

1. **Управление устройствами (Device Management)**  
   Регистрация, жизненный цикл и каталог устройств/датчиков (имя, тип, локация, статус). Граница: всё, что относится к «какое устройство подключено и в каком состоянии оно числится», без сценариев управления и без хранения истории телеметрии.

2. **Управление отоплением (Heating Control)**  
   Команды включения/выключения отопления, управление реле и режимами обогрева. Граница: команды и состояние управления отоплением; опирается на Device Management для идентификации устройств.

3. **Мониторинг температуры / телеметрия (Temperature Monitoring / Telemetry)**  
   Получение, хранение и отдача текущих (и в перспективе исторических) показаний температуры. Граница: данные измерений; не отвечает за команды управления устройствами.

4. **Пользовательский доступ (User Access / Self-Service)** — задел под To-Be  
   Веб-доступ пользователя к функциям управления и мониторинга. В As-Is упрощён (единый веб-интерфейс поверх монолита); в целевой системе — самообслуживание SaaS, выбор модулей и сценарии.

Сейчас все эти контексты **слиты внутри одного монолита** и одной БД `sensors`, что размывает границы и усложняет независимое развитие доменов.

### 4. Проблемы монолитного решения

- **Слабая масштабируемость по частям.** Рост числа домов/датчиков или появление новых типов устройств (свет, ворота, видеонаблюдение) требует масштабировать весь монолит, а не отдельный домен.
- **Синхронная модель end-to-end.** Управление и опрос температуры идут синхронно от сервера к устройству; при росте числа модулей растёт время ответа и риск каскадных задержек/отказов.
- **Связность развёртывания.** Любое изменение (даже в одном домене) требует остановки всего приложения — высокий риск и долгий релизный цикл для команды из 5 разработчиков.
- **Размытые границы доменов.** Device Management, Heating Control и Telemetry живут в одном коде и одной схеме БД — сложно независимо развивать, тестировать и подключать устройства партнёров по стандартным протоколам.
- **Нет самообслуживания.** Подключение выполняет специалист; пользователь не может сам добавить устройство — это блокирует SaaS-модель и масштабирование на несколько регионов.
- **Ограниченный функционал относительно бизнес-целей.** Сейчас только отопление и температура; целевая экосистема требует модульности (свет, ворота, наблюдение, неизвестные типы устройств) — монолит плохо готов к плагинной/модульной модели.

### 5. Визуализация контекста системы — диаграмма C4

Контекстная диаграмма (C4 Context) показывает, как монолитное приложение «Умный дом» взаимодействует с пользователями и модулями отопления/датчиками.

![C4 Context — As-Is «Тёплый дом»](./diagrams/01-c4-context.png)

Исходник PlantUML: [diagrams/01-c4-context.puml](./diagrams/01-c4-context.puml)

Онлайн-просмотр: [PlantUML C4 Context](https://www.plantuml.com/plantuml/uml/dLJBJjjG4DtxAsvP1Ic9DwpeXYX542bjqnA4kee2M10fiP5zez8TJagfb8gegAkgwWjrrOq9k1YS68alcFiB_PBEZ6r8sY2rLQJOcJjpvfmpSxFWoBejswscACxLYfOfZLrPMwlRhMMhxHYrWbDoj7iDSwFPtZJ4jfGxphoks_Nd-Qs6t6wlOuwzGLMcp6zOBNsdMJS9BUV8pMPZNI_EvIeKMdrOrbjrHngsVjidluEPcbOkF7kyMgsjbQhBjVBYqkAZ1vbPJPCDsJG4l8S8Hge3Fbmg5vyZy5HVm1a4uC6fsiFl8NuyaIdEYGHtLlpSUoSAJgxap8lhx_1LlOKh25KFBgWoWk5rg6aLmtOiCqCAic86Fc4yL6_W79y3R48jyVSCdHr1Y7tEA0A-UfM0W9SNS0nNJ6A0oNqyx0Y8i3QYXfW_Gfu-30K4WXCZvLAXwk1xG7suX9iDiSGdB3pigZtL5Ug00pwS2qpruII7DMm67h_4m8LoypFerSeBTBQLQLccrSIn4ELZn1kDAMOT7yXCqWa1Mnc17zCAC1J1ZvIm89l5aeNwFi89wpjXNhqRoilLQYM7PdcfBDNB2imAC3R233V7aASOk2GftLkVuYxn50dqW79uZe5-bta9XuYy4ShbKZmdXLY6FBzn1k6YdIuJGY0SII9e82gM8xTiO-L9USopskAkp6mRTTamjrA7pf1iI0FxZOlgiud73EsnVQaQui5DqirG-yGYu5IKe7fJq58TanPCNuPUOW9l8ChzZ6-dCJwDZzP57UAxIxRoHF51MX0sn9E3nCw0Q17qAxv0OK8HET-VCYFIH3waCxwR6vaMEqyRka_fP0j4h9lYL368wmHSOhE4UqHRZBW3rUMvF3MQoPqSN-UZ6nsyVZgczoP3K-vTzj_l74sOLdj5w6BDM1Uh9UH4b6x9Z6tBbCl_vpp71eU4qXlZgTT_HGQ-J3hALveOweLAoks6qRXZ_xqXsTX94d98YYY9j9mdzwlF8J5nz--cobWqzmuDTy4mD-a__nS0)

# Задание 2. Проектирование микросервисной архитектуры

В этом задании вам нужно предоставить только диаграммы в модели C4. Мы не просим вас отдельно описывать получившиеся микросервисы и то, как вы определили взаимодействия между компонентами To-Be системы. Если вы правильно подготовите диаграммы C4, они и так это покажут.

Декомпозиция To-Be опирается на домены задания 1 (device, heating, telemetry, user) и бизнес-цели кейса (lighting, gate, surveillance, scenario, device-adapter для устройств партнёров и неизвестных типов). Паттерны: Database per Service, API Gateway, синхронный REST для запросов пользователя и асинхронные события через Message Broker для команд и телеметрии.

**Диаграмма контейнеров (Containers)**

![C4 Containers — To-Be «Тёплый дом»](./diagrams/02-c4-containers.png)

Исходник PlantUML: [diagrams/02-c4-containers.puml](./diagrams/02-c4-containers.puml)

**Диаграмма компонентов (Components)**

### device-service

![C4 Components — device-service](./diagrams/03-c4-components-device-service.png)

[diagrams/03-c4-components-device-service.puml](./diagrams/03-c4-components-device-service.puml)

### heating-service

![C4 Components — heating-service](./diagrams/03-c4-components-heating-service.png)

[diagrams/03-c4-components-heating-service.puml](./diagrams/03-c4-components-heating-service.puml)

### lighting-service

![C4 Components — lighting-service](./diagrams/03-c4-components-lighting-service.png)

[diagrams/03-c4-components-lighting-service.puml](./diagrams/03-c4-components-lighting-service.puml)

### gate-service

![C4 Components — gate-service](./diagrams/03-c4-components-gate-service.png)

[diagrams/03-c4-components-gate-service.puml](./diagrams/03-c4-components-gate-service.puml)

### surveillance-service

![C4 Components — surveillance-service](./diagrams/03-c4-components-surveillance-service.png)

[diagrams/03-c4-components-surveillance-service.puml](./diagrams/03-c4-components-surveillance-service.puml)

### telemetry-service

![C4 Components — telemetry-service](./diagrams/03-c4-components-telemetry-service.png)

[diagrams/03-c4-components-telemetry-service.puml](./diagrams/03-c4-components-telemetry-service.puml)

### user-service

![C4 Components — user-service](./diagrams/03-c4-components-user-service.png)

[diagrams/03-c4-components-user-service.puml](./diagrams/03-c4-components-user-service.puml)

### scenario-service

![C4 Components — scenario-service](./diagrams/03-c4-components-scenario-service.png)

[diagrams/03-c4-components-scenario-service.puml](./diagrams/03-c4-components-scenario-service.puml)

### device-adapter

![C4 Components — device-adapter](./diagrams/03-c4-components-device-adapter.png)

[diagrams/03-c4-components-device-adapter.puml](./diagrams/03-c4-components-device-adapter.puml)

**Диаграмма кода (Code)**

### Sequence: пользователь включает отопление

![Code — Sequence heating on](./diagrams/04-code-heating-command-sequence.png)

[diagrams/04-code-heating-command-sequence.puml](./diagrams/04-code-heating-command-sequence.puml)

### Class: критичная модель device-service

![Code — Device classes](./diagrams/04-code-device-classes.png)

[diagrams/04-code-device-classes.puml](./diagrams/04-code-device-classes.puml)

# Задание 3. Разработка ER-диаграммы

### Логическая ER-модель

![ER-диаграмма «Тёплый дом»](./diagrams/05-er-diagram.png)

Исходник PlantUML: [diagrams/05-er-diagram.puml](./diagrams/05-er-diagram.puml)

### Аргументация решений

**1. Логическая модель, а не одна физическая БД.**  
В задании 2 принят паттерн Database per Service. ER здесь — **логическая** модель предметной области: какие сущности и связи есть в экосистеме. Физически они живут в разных БД (`user-db`, `device-db`, `telemetry-db`, `scenario-db` и т.д.). Межсервисные ссылки (например `device_id` в телеметрии) на уровне интеграции — это идентификаторы из событий/API, а не cross-database FK. На диаграмме FK показаны для ясности доменных связей.

**2. Набор сущностей взял из задания + To-Be доменов.**  
Обязательный минимум из условия: `User`, `House`, `Device`, `DeviceType`, `Module`, `TelemetryData`.  
Добавлены `Scenario` / `ScenarioTrigger` / `ScenarioAction`, потому что в бизнес-целях пользователь «настраивает сценарии работы», а в C4 уже есть `scenario-service`. Без этих сущностей модель не покрывала бы To-Be.

**3. `Module` ≠ `Device`.**  
В кейсе модули продаются комплектами и выбираются в self-service, а устройство — конкретный экземпляр в доме.  
- `Module` — каталог комплектов (heating/lighting/gate/…).  
- `Device` — установленный экземпляр с `serial_number`.  
- `HouseModule` (N:M) — какие модули дом подключил/выбрал.  
Так разделяем «что купили/активировали» и «какие физические устройства зарегистрированы».

**4. `DeviceType` — отдельная сущность, не enum в коде.**  
В задании явно: у `Device` есть `type_id` → FK на `DeviceType`. Это нужно для расширения на неизвестные типы партнёров (`protocol`, `UNKNOWN`) без миграции кода при каждом новом типе. Связь `DeviceType` → `Module` (default_module_id) задаёт, к какому комплекту тип обычно относится.

**5. Атрибуты `Device`.**  
По примеру задания: `id`, `type_id`, `house_id`, `serial_number`, `status`.  
Дополнительно: `module_id` (из какого комплекта), `location` (как в монолите `sensors.location`), `partner_protocol_id` (интеграция через `device-adapter`), timestamps.  
`status` хранит жизненный цикл/состояние (`pending|active|offline|on|off|…`) — единая точка для каталога; узкоспециализированное состояние команд (режим отопления, яркость) остаётся в БД control-сервисов и на ER ядра не раздувается.

**6. Связи — строго по формулировкам задания + SaaS.**  

| Связь | Кардинальность | Почему |
| --- | --- | --- |
| User — House | 1:N | «один пользователь — несколько домов; дом связан с одним владельцем» |
| House — Device | 1:N | устройство принадлежит одному дому |
| Device — TelemetryData | 1:N | одно устройство даёт много измерений |
| House — Module | N:M через HouseModule | self-service выбор нескольких комплектов на дом |
| DeviceType — Device | 1:N | классификация экземпляров |
| Module — Device | 1:N | экземпляр порождён комплектом |
| House/User — Scenario | 1:N | сценарии принадлежат дому и автору |
| Scenario — Trigger/Action | 1:N | правило = набор условий и действий над устройствами |

**7. Телеметрия вынесена отдельно от `Device`.**  
В монолите `value` лежало в строке `sensors` — это смешивало каталог и поток измерений. Для масштабирования и истории нужна отдельная `TelemetryData` (`metric_type`, `value`, `unit`, `recorded_at`) с 1:N к устройству. Это соответствует границе `telemetry-service`.

**8. Сценарии нормализованы на Trigger/Action.**  
Условие («температура < 18») и действие («включить отопление») — разные вещи и могут ссылаться на разные устройства. Нормализация упрощает валидацию в `scenario-service` и не тащит гигантский JSON как единственную модель (хотя `action_payload` JSONB оставлен для гибких параметров команды).

**9. Что сознательно не включили в ER.**  
Отдельные таблицы `HeatingState` / `LightingState` / `GateState` / сессии камер — это состояние control-сервисов; на обзорной ER они избыточны и дублировали бы `Device`. Их можно добавить в задании 4/6 как физические схемы сервисов. Общая ER фокусируется на ключевых сущностях экосистемы, как требует задание.

# Задание 4. Создание и документирование API

### 1. Тип API

Для экосистемы «Тёплый дом» выбран **гибридный подход: REST (OpenAPI) + AsyncAPI**.

#### Почему не один тип API

На диаграмме C4 Containers (задание 2) уже зафиксированы два разных стиля связей:

- **синхронный REST** — Web → API Gateway → доменные сервисы (нужен немедленный ответ пользователю);
- **асинхронные Events** — сервисы ↔ Message Broker ↔ device-adapter (команды к железу и поток телеметрии).

Один только REST здесь недостаточен: синхронный опрос сотен/тысяч устройств и ожидание ack по HTTP создают каскадные задержки — это как раз проблема As-Is монолита. Один только AsyncAPI недостаточен для UI: пользователю нужен понятный request/response («зарегистрируй устройство», «какая температура сейчас»).

#### REST API (OpenAPI 3.0) — где и зачем

| Сценарий | Почему REST |
| --- | --- |
| CRUD каталога устройств | Клиенту нужен сразу `201`/`200` с телом ресурса; идемпотентность и стандартные HTTP-коды удобны для Gateway и фронтенда |
| Команда «включить отопление» (вход) | Пользователь ждёт быстрый ответ «команда принята»; долгую доставку на устройство REST не ждёт |
| Чтение последней телеметрии | Классический query: UI/scenario запрашивают актуальное значение и получают JSON |

Контракты: [api/openapi.yaml](api/openapi.yaml)  
Интерактивный просмотр: открыть файл в [Swagger Editor](https://editor.swagger.io/).

#### AsyncAPI — где и зачем

| Сценарий | Почему AsyncAPI / брокер |
| --- | --- |
| `HeatingCommandIssued` | После HTTP `202` команда уходит adapter’у без блокировки пользователя; устройство может ответить с задержкой |
| `TelemetryReceived` | Измерения приходят push-потоком от устройств; telemetry-service и scenario-service подписываются независимо (слабая связанность) |
| `DeviceStatusChanged` | Ack/offline от железа обновляет heating/device без прямых REST-вызовов между ними |
| `DeviceRegistered` / `DeviceUpdated` | Каталог уведомляет подписчиков о новых устройствах без fan-out REST-запросов |

Контракты: [api/asyncapi.yaml](api/asyncapi.yaml)  
Интерактивный просмотр: открыть файл в [AsyncAPI Studio](https://studio.asyncapi.com/).

#### Какие микросервисы покрыты (и почему именно они)

Выбраны три сервиса ядра As-Is → To-Be (отопление + температура + каталог устройств) — они уже есть на ER-диаграмме и на sequence «включить отопление»:

1. **device-service** — источник истины по Device / DeviceType (ER);
2. **heating-service** — домен управления отоплением (главный пользовательский сценарий);
3. **telemetry-service** — хранение и отдача TelemetryData.

Lighting / gate / surveillance используют тот же паттерн «REST-команда → событие → adapter», поэтому отдельно не дублируются в задании (достаточно 4–5 эндпоинтов + async-каналы).

### 2. Документация API

Краткий указатель: [api/README.md](api/README.md)

#### 2.1. REST-эндпоинты (5 штук)

| # | Метод и путь | Сервис | Назначение |
| --- | --- | --- | --- |
| 1 | `POST /api/v1/devices` | device-service | Регистрация устройства в каталоге |
| 2 | `GET /api/v1/devices/{deviceId}` | device-service | Получение информации об устройстве |
| 3 | `PATCH /api/v1/devices/{deviceId}` | device-service | Обновление состояния / атрибутов |
| 4 | `POST /api/v1/heating/{deviceId}/on` | heating-service | Отправка команды «включить отопление» |
| 5 | `GET /api/v1/telemetry/devices/{deviceId}/latest` | telemetry-service | Последнее показание (температура и др.) |

Аргументы выбора эндпоинтов:

- **1–3** закрывают пример из формулировки задания («получение / обновление / команда») для Device Management и согласованы с классовой моделью `DeviceService` (задание 2).
- **4** — точка входа sequence-диаграммы Code: пользовательский синхронный вызов, ответ `202 Accepted`, дальше — брокер.
- **5** — замена синхронного опроса Temperature API в монолите: UI читает уже сохранённую телеметрию, а не блокируется на устройстве.

Полные схемы запросов/ответов, коды `200/201/202/400/404/409/500` и блоки `examples` — в [api/openapi.yaml](api/openapi.yaml).

#### 2.2. Примеры контрактов (кратко)

**Регистрация устройства — запрос**

```json
POST /api/v1/devices
{
  "houseId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "typeId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "moduleId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Отопление гостиная",
  "serialNumber": "HT-2026-00042",
  "location": "Living Room",
  "partnerProtocolId": "mqtt://partner/ht-0042"
}
```

**Ответ `201 Created`**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "houseId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "typeId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "moduleId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Отопление гостиная",
  "serialNumber": "HT-2026-00042",
  "location": "Living Room",
  "status": "pending",
  "partnerProtocolId": "mqtt://partner/ht-0042",
  "createdAt": "2026-11-22T17:30:00Z",
  "updatedAt": "2026-11-22T17:30:00Z"
}
```

**Команда отопления — ответ `202 Accepted`** (исполнение асинхронное)

```json
{
  "commandId": "9b2c1d0e-1111-2222-3333-444455556666",
  "deviceId": "550e8400-e29b-41d4-a716-446655440000",
  "action": "on",
  "mode": "comfort",
  "targetTemperature": 22.5,
  "status": "accepted",
  "acceptedAt": "2026-11-22T18:10:00Z"
}
```

**Последняя телеметрия — ответ `200 OK`**

```json
{
  "id": 10482,
  "deviceId": "550e8400-e29b-41d4-a716-446655440000",
  "metricType": "temperature",
  "value": 21.7,
  "unit": "C",
  "recordedAt": "2026-11-22T18:09:45Z",
  "ingestedAt": "2026-11-22T18:09:46Z"
}
```

#### 2.3. AsyncAPI — каналы и события

| Канал | Событие | Publisher → Consumer | Зачем |
| --- | --- | --- | --- |
| `heating.commands` | `HeatingCommandIssued` | heating-service → device-adapter | Доставка ON/OFF/режима на устройство |
| `telemetry.received` | `TelemetryReceived` | device-adapter → telemetry-service, scenario-service | Поток измерений без синхронного опроса |
| `devices.status` | `DeviceStatusChanged` | device-adapter → device/heating | Факт. статус после ack/timeout |
| `devices.lifecycle` | `DeviceRegistered` / `DeviceUpdated` | device-service → подписчики | Рассылка изменений каталога |

Пример payload `TelemetryReceived`:

```json
{
  "eventId": "c0a80101-0000-4000-8000-000000000002",
  "eventType": "TelemetryReceived",
  "occurredAt": "2026-11-22T18:09:46Z",
  "deviceId": "550e8400-e29b-41d4-a716-446655440000",
  "houseId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "metricType": "temperature",
  "value": 21.7,
  "unit": "C",
  "recordedAt": "2026-11-22T18:09:45Z"
}
```

#### 2.4. Связь REST ↔ Async на сценарии «включить отопление»

1. UI → `POST /api/v1/heating/{deviceId}/on` (REST, OpenAPI).
2. heating-service → `202` + запись в heating-db.
3. Публикация `HeatingCommandIssued` (AsyncAPI).
4. device-adapter → устройство (MQTT/HTTPS).
5. Ack → `DeviceStatusChanged`; измерение → `TelemetryReceived`.
6. UI при необходимости читает `GET .../telemetry/.../latest` (REST).

Так контракты API напрямую продолжают C4 Containers, sequence Code и ER-модель заданий 2–3.

# Задание 5. Работа с docker и docker-compose

### Что сделано

В каталоге `apps/` для задания 5 в `docker-compose.yml` добавлены/настроены:

| Сервис | Назначение | Порт |
| --- | --- | --- |
| `postgres` | PostgreSQL 16 + схема из `./smart_home/init.sql` | 5432 |
| `temperature-api` | Имитация удалённого датчика (рандомная температура) | 8081 |
| `app` (`smart_home`) | Монолит управления датчиками | 8080 |

Код mock-сервиса: `apps/temperature-api/` (Go, Dockerfile multi-stage).

Запуск: `cd apps && ./init.sh` (или `docker compose up --build -d`).

Проверка по заданию (Postman `smarthome-api.postman_collection.json`): **Create Sensor** → **Get All Sensors** несколько раз — поле `value` каждый раз разное.

### Аргументация решений

**1. Язык temperature-api — Go (stdlib `net/http`).**  
Монолит уже на Go, контракт ответа должен совпадать с `TemperatureResponse` в `smart_home/services/temperature_service.go`. Go без внешних зависимостей даёт маленький multi-stage образ и тот же стек, что у ревьюера при разборе `apps/`. Альтернативы (Python/Node) допустимы заданием, но увеличивают разнородность окружения без выигрыша для mock-сервиса.

**2. Два endpoint’а, не один.**  
В задании явно указан `GET /temperature?location=`, но монолит при **Get All Sensors** вызывает `GetTemperatureByID` → `GET /temperature/{sensorID}`. Без path-роута проверка «Create + Get All» не покажет живую температуру. Поэтому реализованы оба варианта; маппинг location ↔ sensorId взят из шаблона задания (1=Living Room, 2=Bedroom, 3=Kitchen).

**3. Формат JSON ответа.**  
Поля `value`, `unit`, `timestamp`, `location`, `status`, `sensor_id`, `sensor_type`, `description` совпадают со структурой клиента в монолите. `status: "active"` — чтобы UI/ответ выглядели как рабочий датчик; `value` генерируется заново на каждый запрос в диапазоне ~15–30°C.

**4. Порт 8081 и имя сервиса `temperature-api`.**  
Требование задания + уже заданный в шаблоне `TEMPERATURE_API_URL=http://temperature-api:8081`. DNS-имя сервиса в Docker Compose-сети совпадает с hostname в URL монолита — отдельная конфигурация hosts не нужна.

**5. Postgres: credentials и init.sql.**  
`DATABASE_URL` монолита: `postgres://postgres:postgres@postgres:5432/smarthome`. Соответственно `POSTGRES_USER` / `POSTGRES_PASSWORD`. Скрипт `./smart_home/init.sql` монтируется в `/docker-entrypoint-initdb.d/` — стандартный механизм официального образа Postgres.  
**Почему не задаём `POSTGRES_DB=smarthome`:** `init.sql` сам делает `CREATE DATABASE smarthome` и `\c smarthome`. Если БД уже создана entrypoint’ом, `CREATE` падает при `ON_ERROR_STOP=1`, и таблица `sensors` не создаётся. Поэтому БД создаёт только init-скрипт.

**6. Healthcheck Postgres + `depends_on`.**  
У `app` уже было `condition: service_healthy` для Postgres — добавлен `healthcheck` (`pg_isready`), иначе Compose не дождётся готовности БД и монолит упадёт на старте. Для `temperature-api` достаточно `service_started`: это stateless HTTP без миграций.

**7. Multi-stage Dockerfile для temperature-api.**  
Как у `smart_home`: сборка в `golang:1.22-alpine`, runtime в `alpine` — меньше атакующая поверхность и размер образа, бинарь статически слинкован (`CGO_ENABLED=0`).


# **Задание 6. Разработка MVP**

Необходимо создать новые микросервисы и обеспечить их интеграции с существующим монолитом для плавного перехода к микросервисной архитектуре.

### Что сделано

Реализован MVP перехода к микросервисам по паттерну **Strangler Fig**:

| Компонент | Технология | Порт | Роль |
| --- | --- | --- | --- |
| `smart_home` (монолит) | Go | 8080 | Источник истины для `/api/v1/sensors`, публикует события |
| `temperature-api` | Go | 8081 | Внешний источник температуры (задание 5) |
| `device-service` | **Java 17 / Spring Boot** | 8082 | Каталог устройств (Device Management) |
| `telemetry-service` | **Python / FastAPI** | 8083 | Хранение показаний (Telemetry) |
| `rabbitmq` | RabbitMQ 3.13 | 5672 / 15672 | Шина событий |
| `postgres` | PostgreSQL 16 | 5432 | `smarthome` + `device_db` + `telemetry_db` |

Код: `apps/device-service`, `apps/telemetry-service`, интеграция в `apps/smart_home/messaging`.  
Запуск: `cd apps && docker compose up --build -d` (или `./init.sh`).

### Аргументация решений

**1. Какие сервисы вынесены первыми и почему**  
В задании 1/2 границы доменов уже зафиксированы: Device Management и Telemetry — самые «чистые» кандидаты на вынос. У них слабая связанность с командами отопления и понятные CRUD/append-only сценарии. Heating/lighting/gate оставили в To-Be: для MVP достаточно доказать миграционный контур, а не клонировать всю экосистему.

**2. Почему разные ООП-языки**  
Требование: «каждый микросервис на своем ООП языке».  
- `device-service` — **Java**: классический OOP (инкапсуляция `Device`, application-сервис, JPA-репозиторий), совпадает с вариантом `Go/Java` на C4 Containers.  
- `telemetry-service` — **Python**: OOP через классы `Measurement`, `TelemetryRepository`, `TelemetryService`, `TelemetryEventConsumer`. Язык другой, модель та же (слои API → application → domain → persistence).  
Монолит остаётся на Go: его не переписываем, а «обволакиваем» событиями.

**3. Почему RabbitMQ, а не только REST**  
На C4 Containers уже заложена шина. Для плавного перехода синхронный REST монолит→сервис создавал бы жёсткую связность и каскадные отказы. Topic-exchange `warmhouse.events` даёт:
- асинхронную репликацию каталога (`device.registered|updated|removed`);
- поток измерений (`telemetry.received`);
- независимое масштабирование consumer’ов;
- устойчивость: если микросервис недоступен, монолит продолжает отвечать клиенту (событие будет обработано после восстановления / при необходимости можно добавить outbox позже).

**4. Паттерн интеграции — Strangler Fig**  
Клиент по-прежнему ходит в монолит (`POST/GET /api/v1/sensors`). После успешных мутаций и при чтении температуры монолит **публикует события**. Микросервисы строят свои проекции и отдают их по своим REST API:
- `GET/POST /api/v1/devices` → device-service  
- `GET/POST /api/v1/telemetry` → telemetry-service  

Так функциональность постепенно «переезжает» без big-bang rewrite. Следующий шаг To-Be — API Gateway начнёт маршрутизировать часть запросов напрямую в сервисы.

**5. Database per Service**  
На одном инстансе Postgres созданы отдельные БД `device_db` и `telemetry_db` (`init-service-dbs.sql` при первом init тома + one-shot сервис `db-init` / `ensure-service-dbs.sh` для уже существующих томов). Это компромисс MVP: физическая изоляция схем как в задании 3, без трёх отдельных контейнеров Postgres. Cross-DB FK нет: `device_id` в телеметрии — интеграционный идентификатор из события.

**6. Простейшая логика сервисов (осознанно)**  
- device-service: upsert каталога по событиям + CRUD для проверки API.  
- telemetry-service: append-only запись измерений + выборка истории/latest.  
Без auth, без саг и без полной ER из задания 3 — цель MVP: доказать границы, контракты и docker-запуск.

**7. Контракт событий**

```json
// device.*
{ "deviceId": 1, "name": "...", "type": "temperature", "location": "...", "status": "active", "eventType": "DeviceRegistered" }

// telemetry.received
{ "deviceId": 1, "metricType": "temperature", "value": 22.5, "unit": "C", "location": "...", "status": "active", "recordedAt": "...", "eventType": "TelemetryReceived" }
```

**8. Как проверить**

```bash
cd apps
docker compose up --build -d

# 1) создать датчик в монолите → событие DeviceRegistered
curl -s -X POST http://localhost:8080/api/v1/sensors \
  -H 'Content-Type: application/json' \
  -d '{"name":"Living Room Sensor","type":"temperature","location":"Living Room","unit":"C"}'

# 2) каталог в device-service
curl -s http://localhost:8082/api/v1/devices

# 3) чтение сенсоров → TelemetryReceived
curl -s http://localhost:8080/api/v1/sensors
curl -s http://localhost:8083/api/v1/telemetry
```

RabbitMQ UI: http://localhost:15672 (guest/guest).
