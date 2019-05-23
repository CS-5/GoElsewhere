import React from "react";
import Card from "@material-ui/core/Card";
import CardHeader from "@material-ui/core/CardHeader";
import CardContent from "@material-ui/core/CardContent";
import CardActions from "@material-ui/core/CardActions";
import IconButton from "@material-ui/core/IconButton";
import Microlink from "@microlink/react";
import { makeStyles } from "@material-ui/core/styles";
import LinkIcon from "@material-ui/icons/Link";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import Tooltip from "@material-ui/core/Tooltip";
import Zoom from "@material-ui/core/Zoom";
import Typeography from "@material-ui/core/Typography";
import copy from "copy-to-clipboard";

const useStyles = makeStyles(theme => {
  return {
    cardHeader: {
      backgroundColor: theme.palette.primary.main,
      height: 75
    },
    cardAvatar: {
      fontSize: "large"
    },
    avatarButton: {
      padding: 5,
      color: "white",
      "&:hover": {
        color: theme.palette.secondary.main
      }
    },
    cardTitle: {
      color: "white",
      fontFamily: "Roboto",
      fontWeight: 200,
      fontSize: "1.75em"
    },
    cardSubheader: {
      color: "white"
    },

    cardActionsButtons: {
      marginLeft: "auto"
    },
    button: {
      color: theme.palette.primary.main,
      "&:hover": {
        color: theme.palette.secondary.main
      }
    }
  };
});

export default function RedirectCard({
  code = 111111,
  link = "/" + code,
  url = "https://google.com",
  created = "1/11/2000",
  hits = 0,
  deleter
}) {
  const classes = useStyles();

  function clip() {
    copy(link);
  }

  function del() {
    deleter(code);
  }

  return (
    <Card>
      <CardHeader
        classes={{
          root: classes.cardHeader,
          avatar: classes.cardAvatar,
          title: classes.cardTitle,
          subheader: classes.cardSubheader
        }}
        avatar={
          <Tooltip
            title="Copy Link"
            aria-label="Copy Link"
            TransitionComponent={Zoom}
          >
            <IconButton classes={{ root: classes.avatarButton }} onClick={clip}>
              <LinkIcon fontSize={"large"} />
            </IconButton>
          </Tooltip>
        }
        title={code}
        subheader={<>Created: {created}</>}
      />

      <CardContent>
        <Microlink url={url} media={["screenshot", "image", "logo"]} />
      </CardContent>
      <CardActions disableSpacing>
        <Typeography>Hits: {hits}</Typeography>
        <div className={classes.cardActionsButtons}>
          <Tooltip
            title="Edit (WIP)"
            aria-label="Edit (WIP)"
            TransitionComponent={Zoom}
          >
            <IconButton
              size="small"
              variant="contained"
              classes={{ root: classes.button }}
            >
              <EditIcon />
            </IconButton>
          </Tooltip>
          <Tooltip
            title="Delete"
            aria-label="Delete"
            TransitionComponent={Zoom}
            onClick={del}
          >
            <IconButton
              size="small"
              variant="contained"
              classes={{ root: classes.button }}
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </div>
      </CardActions>
    </Card>
  );
}
