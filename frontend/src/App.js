import React, { Component } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

import { Table, Container, Row, Col } from 'reactstrap';
const API_BASE_URL = process.env.API_BASE_URL || "http://localhost:8080"
class App extends Component {
    state = {
        data: [],
    };
    constructor(props) {
        super(props);
        this.fetchData = this.fetchData.bind(this);
    }
    componentDidMount() {
        this.fetchData();
    }
    fetchData() {
        fetch(API_BASE_URL + "/users")
            .then(response => response.json())
            .then(data => this.setState({ data: data.data }));
    }
  render() {
      const items  = this.state.data;
    return (
        <Container>
            <Row>
                <Col>
                  <div className="App">
                    <header className="App-header">
                      <img src="https://s3-eu-west-1.amazonaws.com/ddme-datixuk/theme/logos/datix-logo-wo.svg" className="App-logo" alt="logo" />
                    </header>
                  </div>
                </Col>
            </Row>
              <Row>
                  <Col>
                      <Table striped bordered hover>
                          <thead>
                          <tr>
                              <th>#</th>
                              <th>First Name</th>
                              <th>Last Name</th>
                              <th>Email</th>
                              <th>Gender</th>
                              <th>IP Address</th>
                          </tr>
                          </thead>
                          <tbody>
                          {items.map((row, key) =>
                          <tr key={key}>
                              <th scope="row">{row.Id}</th>
                              <td>{row.first_name}</td>
                              <td>{row.last_name}</td>
                              <td>{row.email}</td>
                              <td>{row.gender}</td>
                              <td>{row.ip_address}</td>
                          </tr>
                          )}
                          </tbody>
                      </Table>

                  </Col>
              </Row>
          </Container>
    );
  }
}

export default App;
