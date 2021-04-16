import axios from "axios";
import config from "./config";

export const call = (opts) => {

  return new Promise((resolve, reject) => {

    let headers = opts.headers || {};

    let token = localStorage.getItem("token");
    if (token) {
      headers["authorization"] = token;
    }

    headers["content-type"] = "application/json";

    axios({
      url: process.env.VUE_APP_ROOT_API + opts.path,
      method: opts.method.toUpperCase(),
      headers: headers,
      data: opts.body,
      params: opts.params || {}
    })
      .then(res => {
        resolve(res);
      })
      .catch(err => {
        reject(err);
      });

  });

};

export const raw = async (opts) => {
  try {
    let res = await axios(opts);

    return res;
  }
  catch (err) {
    throw new Error(err);
  }
};
