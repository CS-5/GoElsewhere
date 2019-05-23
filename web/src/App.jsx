import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import { ThemeProvider } from "@material-ui/styles";
import { createMuiTheme } from "@material-ui/core/styles";
import RedirectList from "./components/RedirectList";
import Nav from "./components/Nav";
import CreateDialog from "./components/CreateDialog";
import ErrorDialog from "./components/ErrorDialog";
import { CssBaseline } from "@material-ui/core";
import axios from "axios";

function Login() {
  return <h1>Login</h1>;
}

export default function App() {
  const [redirects, setRedirects] = useState([]);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [errorDialogOpen, setErrorDialogOpen] = useState(false);
  const [errorDialogMessage, setErrorDialogMessage] = useState("");
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const theme = createMuiTheme({
    palette: {
      primary: {
        main: "#0d47a1"
      },
      secondary: {
        main: "#e65100"
      }
    },
    overrides: {
      MuiTooltip: {
        tooltip: {
          fontSize: "0.8em"
        }
      }
    }
  });

  useEffect(() => {
    getRedirects().then(reds => setRedirects(reds));
  }, []);

  function openCreateDialog() {
    setCreateDialogOpen(true);
  }

  function closeCreateDialog() {
    setCreateDialogOpen(false);
  }

  function closeErrorDialog() {
    setErrorDialogOpen(false);
  }

  function showError(msg) {
    setErrorDialogMessage(msg);
    setErrorDialogOpen(true);
  }

  function getRedirects() {
    return axios
      .get("/api/list")
      .then(response => {
        let r = [];
        const data = response.data;

        for (var key in data) {
          if (data.hasOwnProperty(key)) {
            r.push({
              id: data[key].id,
              code: key,
              url: data[key].url,
              link: data[key].link,
              created: data[key].created,
              hits: data[key].hits
            });
          }
        }

        return r;
      })
      .catch(function(error) {
        console.log(error);
      });
  }

  function createRedirect(data) {
    fetch("/api/create", {
      method: "POST",
      body: data
    })
      .then(data => data.json())
      .then(data => {
        if (data.good) {
          let entry = data.entry;

          let r = {
            id: entry.id,
            code: entry.code,
            url: entry.url,
            link: entry.link,
            created: entry.created,
            hits: entry.hits
          };

          setRedirects([...redirects, r]);
        } else {
          showError(data.error);
        }
      });
  }

  function deleteRedirect(c) {
    fetch("/api/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code: c })
    })
      .then(data => data.json())
      .then(data => {
        if (data.good) {
          let newRedirects = [...redirects];
          newRedirects.splice(redirects.findIndex(i => i.code === c), 1);

          setRedirects(newRedirects);
        } else {
          showError(data.error);
        }
      });
  }

  return (
    <Router basename={"/admin"}>
      <CssBaseline />
      <ThemeProvider theme={theme}>
        <Nav createOpener={openCreateDialog} />
        <CreateDialog
          open={createDialogOpen}
          closer={closeCreateDialog}
          creater={createRedirect}
        />
        <ErrorDialog
          open={errorDialogOpen}
          closer={closeErrorDialog}
          message={errorDialogMessage}
        />

        <div style={{ padding: 10, paddingTop: 74 }}>
          <Route
            exact
            path="/"
            render={props => (
              <RedirectList
                redirects={redirects}
                createDialog={openCreateDialog}
                deleteDialog={deleteRedirect}
              />
            )}
          />
          <Route path="/login" component={Login} />
        </div>
      </ThemeProvider>
    </Router>
  );
}
