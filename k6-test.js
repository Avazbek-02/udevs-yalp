import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: "30", target: 3 },
        { duration: "20", target: 100 },
        { duration: "300", target: 50 },
    ],
    cloud: {
        projectID: "3739096",
        name: "First Project"
    }
};

export default () => {
    let Data_Login = JSON.stringify({
        email: "avazbekmambetov9@gmail.com",
        password: "1234",
        platform: "web",
    });

    let uniqueId = Math.random().toString(36).substring(2, 8);
    let emailDomain = "gmail.com";

    let Data_Register = JSON.stringify({
        full_name: "Test",
        user_type: "user",
        user_role: "user",
        username: "testusername",
        email: `test+${uniqueId}@${emailDomain}`,
        profile_picture: `${uniqueId}`,
        status: "inverify",
        password: "1234",
        gender: "male",
    });
    

    let loginParams = {
        headers: {
            "Content-Type": "application/json"
        }
    };

    // Login request
    const resLogin = http.post('http://localhost:8080/v1/auth/login', Data_Login, loginParams);
    let token = resLogin.json('user.access_token');
    let userid = resLogin.json('user.id');
    console.log('Access Token:', token);

    let registerParams = {
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        }
    };


    check(resLogin, {
        "status code 200": (r) => r.status === 200
    });

    const resRegister = http.post('http://localhost:8080/v1/user/', Data_Register, registerParams);

    check(resRegister, {
        "status code 201": (r) => r.status === 201
    });

    const resGetSingleUser = http.get(`http://localhost:8080/v1/user/${userid}`, registerParams);

    check(resGetSingleUser, {
        "status code 200": (r) => r.status === 200
    });

    console.log('Register Response:', resRegister.body);
    console.log('Get Single User Response:', resGetSingleUser.body);
};
