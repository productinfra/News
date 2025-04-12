/**
 *  @param {String/Number} date
 *  @param {format}
 *
 */
function formatDate(date, format) {
  format = format || "";
  let time = new Date(date);
  let year = time.getFullYear();
  let month = time.getMonth() + 1;
  let day = time.getDate();
  let hour = time.getHours();
  let minute = time.getMinutes();
  let second = time.getSeconds();
  let milliSecond = time.getMilliseconds();
  let resutTimeObject = {
    year,
    month,
    day,
    hour,
    minute,
    second,
    milliSecond,
  };
  if (!format) {
    return resutTimeObject;
  } else {
    format = format.trim();
    if (format === "YYYY-MM-DD hh:mm:ss") {
      return (
        year +
        "-" +
        month +
        "-" +
        day +
        " " +
        hour +
        ":" +
        minute +
        ":" +
        second
      );
    } else if (format === "YYYY/MM/DD hh:mm:ss") {
      return (
        year +
        "/" +
        month +
        "/" +
        day +
        " " +
        hour +
        ":" +
        minute +
        ":" +
        second
      );
    } else if (format === "YYYY-MM-DD") {
      return year + "-" + month + "-" + day;
    } else if (format === "YYYY/MM/DD") {
      return year + "/" + month + "/" + day;
    } else if (format === "hh:mm:ss") {
      return hour + ":" + minute + ":" + second;
    } else if (format === "hh-mm-ss") {
      return hour + "-" + minute + "-" + second;
    } else {
      return date;
    }
  }
}

/**
 *
 *  @param {Number} num
 *  @param {Number} n
 *  @param {Number} num1
 *  @returns {String}
 */

function fixedNumber(num, n, num1) {
  num1 = num1 || 0;
  return (Array(n).join(num1) + num).slice(-n);
}

/**
 *
 * @param {Number} unitSecond
 * @param {String} format
 * @returns
 */
function handleTime(unitSecond, format) {
  let timeArr = new Array(0, 0, 0, 0, 0);
  let unitYear = 365 * 24 * 60 * 60;
  let unitDay = 24 * 60 * 60;
  let unitHour = 60 * 60;
  let unitMin = 60;
  let unitSec = 0;
  if (!unitSecond) {
    return;
  }
  if ((format ? format.indexOf("Y") > -1 : false) && unitSecond >= unitYear) {
    timeArr[0] = parseInt(unitSecond / unitYear);
    unitSecond %= unitYear;
  }
  if (unitSecond >= unitDay) {
    timeArr[1] = parseInt(unitSecond / unitDay);
    unitSecond %= unitDay;
  }
  if (unitSecond >= unitHour) {
    timeArr[2] = parseInt(unitSecond / unitHour);
    unitSecond %= unitHour;
  }
  if (unitSecond >= unitMin) {
    timeArr[3] = parseInt(unitSecond / unitMin);
    unitSecond %= unitMin;
  }
  if (unitSecond > unitSec) {
    timeArr[4] = unitSecond;
  }
  return timeArr;
}

/**
 * @param {Number} year
 * @param {Number} month
 * @param {Number} day
 * @param {Number} hour
 * @param {Number} minute
 * @param {Number} second
 * @returns
 */
function getTime(year, month, day, hour, minute, second) {
  let startTime = Math.round(
    new Date(Date.UTC(year, month - 1, day, hour, minute, second)).getTime() /
      1000
  );
  let nowTime = Math.round((new Date().getTime() + 8 * 60 * 60 * 1000) / 1000);
  return handleTime(nowTime - startTime);
}

/**
 * @param {Array} runTimeArr
 * @param {String} format
 * @returns
 */
function handleResult(runTimeArr, format) {
  if (!format) {
    return (
      runTimeArr[1] +
      " Days " +
      runTimeArr[2] +
      " Hours " +
      runTimeArr[3] +
      " Minutes " +
      runTimeArr[4] +
      " Seconds"
    );
  } else if (format === "Y-D h:m:s") {
    return (
      runTimeArr[0] +
      "Years" +
      runTimeArr[1] +
      "Day" +
      runTimeArr[2] +
      "Hours" +
      runTimeArr[3] +
      "Minutes" +
      runTimeArr[4] +
      "Seconds"
    );
  } else if (format === "D h:m:s") {
    return (
      runTimeArr[1] +
      "Day" +
      runTimeArr[2] +
      "Hours" +
      runTimeArr[3] +
      "Minutes" +
      runTimeArr[4] +
      "Seconds"
    );
  }
}

/**
 * @param {String/Number/Object} timeStamp
 * @param {String/DOMObject} el
 * @param {String} desc
 * @param {Number} year
 * @param {Number} month
 * @param {Number} day
 * @param {Number} hour
 * @param {Number} minute
 * @param {Number} second
 * @param {boolean} flag
 */
export function runTime({
  el,
  timeStamp,
  desc,
  year,
  month,
  day,
  hour,
  minute,
  second,
  flag = true,
  format,
}) {
  desc = desc || "";
  if (timeStamp) {
    let time = formatDate(timeStamp);
    year = time.year;
    month = time.month;
    day = time.day;
    hour = time.hour;
    minute = time.minute;
    second = time.second;
  }
  if (flag && el) {
    let time_wrapper = document.querySelector(el);
    setInterval(() => {
      let runTimeArr = getTime(year, month, day, hour, minute, second);
      runTimeArr[2] = fixedNumber(runTimeArr[2], 2, 0);
      runTimeArr[3] = fixedNumber(runTimeArr[3], 2, 0);
      runTimeArr[4] = fixedNumber(runTimeArr[4], 2, 0);
      time_wrapper.innerText = desc + handleResult(runTimeArr, format);
    }, 1000);
  } else {
    return getTime(year, month, day, hour, minute, second);
  }
}
