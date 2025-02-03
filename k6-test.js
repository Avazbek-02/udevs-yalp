import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: "5s", target: 1000 },
        { duration: "30s", target: 1000 },
        { duration: "5s", target: 0 },
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

    // let uniqueId = Math.random().toString(36).substring(2, 8);
    // let emailDomain = "gmail.com";

    // let Data_Register = JSON.stringify({
    //     full_name: "Test",
    //     user_type: "user",
    //     user_role: "user",
    //     username: "testusername",
    //     email: `test+${uniqueId}@${emailDomain}`,
    //     profile_picture: `${uniqueId}`,
    //     status: "inverify",
    //     password: "1234",
    //     gender: "male",
    // });
    

    // let loginParams = {
    //     headers: {
    //         "Content-Type": "application/json"
    //     }
    // };

    

    let registerParams = {
        headers: {
            "Authorization": `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwbGF0Zm9ybSI6IndlYiIsInNlc3Npb25faWQiOiI4NjZlZGNlYS0zYzkxLTQwNjctOTY1NS1lYzM5MmQyYzZhMTQiLCJzdWIiOiI1MDM1OTQ0MS04YjdkLTQ1ZWItYTgzZC04MjMxMmJlOGM5NWQiLCJ1c2VyX3JvbGUiOiJ1c2VyIiwidXNlcl90eXBlIjoidXNlciJ9.7mOZOXkDMX1qX9jglq4ZDcGm_35FqQuzVW0BAfYk9qM`,
            "Content-Type": "application/json",
        }
    };


  

    // const resRegister = http.post('http://localhost:8080/v1/user/', Data_Register, registerParams);

    // check(resRegister, {
    //     "status code 201": (r) => r.status === 201
    // });

    const resGetSingleUser = http.get(`http://localhost:8080/v1/user/50359441-8b7d-45eb-a83d-82312be8c95d`, registerParams);

    check(resGetSingleUser, {
        "status code 200": (r) => r.status === 200
    });

    // console.log('Register Response:', resRegister.body);
    // console.log('Get Single User Response:', resGetSingleUser.body);
};

