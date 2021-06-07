package utils

const ONLY_FILTERED_LANGUAGES = "only_lang"
const SchoolIdValidateRegex = `^[0-9a-z_\-\.]{3,64}$`
const CourseIdValidateRegex = `^[0-9a-zA-Z_\-\.]{3,64}$`

const TRANSPRO_API_BASE_URL string = "https://bcause-api.com"
const TRANSPRO_API_LOGIN_URL string = TRANSPRO_API_BASE_URL + "/translator/login"

const TRANSPRO_API_AVAILABLE_PROJECTLIST_URL = TRANSPRO_API_BASE_URL + "/order/getAvailableWebOrders"

const QUICKTRANSLATE_ROOT_URL string = "https://quicktranslate.com"
const QUICKTRANSLATE_LOGIN_URL string = QUICKTRANSLATE_ROOT_URL + "/login?_url=/"
const QUICKTRANSLATE_JOB_FETCH_URL string = QUICKTRANSLATE_ROOT_URL + "/api/translator/jobs/web_api"

// const JOB_DELIVERED_LIST string = quickTranslateRootUrl + "/api/translator/jobs/web_api/delivered"

const QUICKTRANSLATE_JOBAPI_ACCEPT_URL string = QUICKTRANSLATE_ROOT_URL + "/api/translator/web/job/accept"
const QUICKTRANSLATE_LOGIN_API_URL string = QUICKTRANSLATE_ROOT_URL + "/auth/login"
const QUICKTRANSLATE_JOB_DETAIL_URL string = QUICKTRANSLATE_ROOT_URL + "/mypage_trans/job/%v"

const STEPES_TRANSLATOR_ROOT_URL string = "https://translator.stepes.com"
const STEPES_TRANSLATOR_LOGIN_URL string = STEPES_TRANSLATOR_ROOT_URL + "/translator?returnto="
const STEPES_CHECKJOB_URL string = STEPES_TRANSLATOR_ROOT_URL + "/stepes-check-jobs.html"
