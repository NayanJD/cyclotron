import grpc from 'k6/net/grpc';
import { check, sleep } from "k6";
import { Faker } from "k6/x/faker";

export const options = {
  vus: 1000,
  // duration: '50s',
  stages: [
    { duration: '5s', target: 100 },
    { duration: '5s', target: 200 },
    { duration: '5s', target: 300 },
    { duration: '5s', target: 400 },
    { duration: '5s', target: 500 },
    { duration: '5s', target: 600 },
    { duration: '5s', target: 700 },
    { duration: '5s', target: 800 },
    { duration: '5s', target: 900 },
    { duration: '5s', target: 1000 },

  ],
};

const client = new grpc.Client();

client.load(['../..'], 'pkg/grpc/pb/user.proto');

export default () => {
  client.connect(__ENV.USER_SVC_URL + ':8082', {
    plaintext: true
  });
    
  let f = new Faker();
    let firstName = f.firstName()
    let lastName = f.lastName()
    let username = f.email();
    let password = f.password(true, true, true, true, false, 12)
    
  let data = {
      Username: username,
      FirstName: firstName,
      LastName: lastName,
      Password: password,
      Dob: "2000-01-27T15:07:56.572316900Z"
      // Dob: {
      //   seconds: "966757140",
      //   nanos: 0
      // }
    }

  let response = client.invoke('pb.User/Register', data);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
      });

  data = { Username: username, Password: password };
  response = client.invoke('pb.User/Login', data);

  console.log("login response", response)
  let accessToken = response.message["AccessToken"]
  let refreshToken = response.message["RefreshToken"]

  check(response, {
    'status is OK': (r) => r && r.status === grpc.StatusOK,
  });

  let count = 0
  while(count < 13) {
    console.log("refreshToken", refreshToken)
    data = {
        RefreshToken: refreshToken
    }

    response = client.invoke('pb.User/RefreshAccessToken', data);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
      });

    accessToken = response.message["AccessToken"]
    refreshToken = response.message["RefreshToken"]

    sleep(2)
    count++;
  }

  // console.log(JSON.stringify(response.message));

  client.close();
  // sleep(1);
};
