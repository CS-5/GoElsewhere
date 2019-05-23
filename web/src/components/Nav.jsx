import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";

const useStyles = makeStyles(theme => ({
  grow: {
    flexGrow: 1
  },
  menuTitle: {
    marginRight: 20
  },
  menuPaper: {
    width: 100
  }
}));

export default function Nav({ createOpener }) {
  const classes = useStyles();

  return (
    <div className={classes.grow}>
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" className={classes.menuTitle}>
            Go-Elsewhere
          </Typography>
          <div className={classes.grow} />
          <Button color="inherit">Logout</Button>
        </Toolbar>
      </AppBar>
    </div>
  );
}
